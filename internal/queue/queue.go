package queue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/afiskon/promtail-client/promtail"
	"github.com/go-kit/kit/sd/lb"
	"github.com/nats-io/stan.go"
	"github.com/rusrafkasimov/catalogs/internal/config"
	"github.com/rusrafkasimov/catalogs/internal/trace"
	"github.com/rusrafkasimov/catalogs/pkg/models"
	"sync"
	"time"
)

var (
	errNoConnection = errors.New("no connection to NATS system")
	errQueueClosed  = errors.New("queue already closed")
)

type Event interface {
	Operation() *models.Operation
	Sequence() uint64
	Ack() error
}

type event struct {
	opn *models.Operation
	seq uint64
	ack func() error
}


func (e *event) Operation() *models.Operation {
	return e.opn
}

func (e *event) Sequence() uint64 {
	return e.seq
}

func (e *event) Ack() error {
	return e.ack()
}


// Queue is a NATS-based event queue. It allows subscribing to events or publish them.

type EventQueue interface {
	Publish(op *models.Operation) error
	Subscribe() (<-chan Event, error)
}

type Queue struct {
	writeOnly        bool
	logger           promtail.Client
	nodeID           string
	url              string
	ackWait          time.Duration
	reconnectTimeout time.Duration
	clusterID        string
	subject          string

	mu             sync.RWMutex
	conn           stan.Conn
	input          chan Event
	output         chan Event
	sequenceNumber uint64
	closed         bool

	wg     sync.WaitGroup
	doneCh chan struct{}

	now func() time.Time
}


// NewQueue creates a new Queue.
func NewQueue(ctx context.Context, logger promtail.Client, cfg *config.Configuration) (*Queue, error) {
	URL, err := cfg.Get("EVENT_QUEUE_URL")
	if err != nil {
		fmt.Println(err)
	}

	ClusterID, err := cfg.Get("EVENT_QUEUE_CLUSTER_ID")
	if err != nil {
		fmt.Println(err)
	}

	Subject, err := cfg.Get("EVENT_QUEUE_SUBJECT")
	if err != nil {
		fmt.Println(err)
	}

	clientId := ctx.Value("Name")
	q := &Queue{
		writeOnly:        false,
		logger:           logger,
		nodeID:           clientId.(string),
		url:              URL,
		ackWait:          time.Second,
		reconnectTimeout: time.Second,
		clusterID:        ClusterID,
		subject:          Subject,
		mu:               sync.RWMutex{},
		conn:             nil,
		input:            make(chan Event),
		output:           make(chan Event),
		sequenceNumber:   0,
		closed:           false,
		wg:               sync.WaitGroup{},
		doneCh:           make(chan struct{}),
		now:              time.Now,
	}

	err = q.connect()
	if err != nil {
		trace.OnError(logger, nil, err)
		q.dialBackground()
	}

	if !q.writeOnly {
		q.handleEvents()
	}

	return q, nil
}


// connect establish new connection.
func (q *Queue) connect() error {
	q.mu.Lock()
	defer q.mu.Unlock()

	stanOpts := []stan.Option{
		stan.SetConnectionLostHandler(q.recoverConn),
		stan.NatsURL(q.url),
	}

	conn, err := stan.Connect(q.clusterID, q.nodeID, stanOpts...)
	if err != nil {
		return err
	}

	q.logger.Infof("Established NATS Streaming connection")

	if !q.writeOnly {
		subsOpts := []stan.SubscriptionOption{
			stan.MaxInflight(1),
			stan.AckWait(q.ackWait),
			stan.SetManualAckMode(),
		}

		if q.sequenceNumber > 0 {
			subsOpts = append(subsOpts, stan.StartAtSequence(q.sequenceNumber))
		} else {
			subsOpts = append(subsOpts, stan.DeliverAllAvailable())
		}

		_, err = conn.QueueSubscribe(q.subject, q.nodeID, q.handleMessage, subsOpts...)
		if err != nil {
			_ = conn.Close()
			return fmt.Errorf("failed to subscribe to queue: %w", err)
		}
	}

	q.conn = conn

	return nil
}

// getConn return connection with mutex.
func (q *Queue) getConn() (stan.Conn, error) {
	q.mu.RLock()
	defer q.mu.RUnlock()

	if q.conn == nil {
		return nil, errNoConnection
	}

	return q.conn, nil
}

// recoverConn notifies about disconnection and start dialBackground.
func (q *Queue) recoverConn(_ stan.Conn, reason error) {
	q.logger.Errorf("NATS connection lost: ", reason)
	q.dialBackground()
}

// dialBackground tries to reestablish connection by timer.
func (q *Queue) dialBackground() {
	q.wg.Add(1)
	go func() {
		tc := time.NewTicker(q.reconnectTimeout)
		defer func() {
			tc.Stop()
			q.wg.Done()
		}()

		for {
			select {
			case <-tc.C:
				err := q.connect()
				if err != nil {
					trace.OnError(q.logger, nil, err)
					continue
				}
				return

			case <-q.doneCh:
				return
			}
		}
	}()
}


// Subscribe returns channel with replication events.
func (q *Queue) Subscribe() (<-chan Event, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.closed {
		return nil, errQueueClosed
	}

	return q.output, nil
}

// Publish sends operation in replication queue.
func (q *Queue) Publish(op *models.Operation) error {
	conn, err := q.getConn()
	if err != nil {
		return err
	}

	op.Timestamp = q.now()
	data, err := json.Marshal(op)
	if err != nil {
		return fmt.Errorf("failed to encode data: %w", err)
	}

	return conn.Publish(q.subject, data)
}


// handleMessage get and unmarshal message, send event to input chan.
func (q *Queue) handleMessage(msg *stan.Msg) {
	op := &models.Operation{}
	err := json.Unmarshal(msg.Data, op)
	if err != nil {
		trace.OnError(q.logger, nil, err)
		return
	}

	q.mu.Lock()
	q.sequenceNumber = msg.Sequence
	q.mu.Unlock()

	event := &event{
		opn: op,
		seq: msg.Sequence,
		ack: msg.Ack,
	}

	select {
	case q.input <- event:
	case <-q.doneCh:
	}
}

// handleEvents handle events and send them to output chan.
func (q *Queue) handleEvents() {
	q.wg.Add(1)
	go func() {
		defer q.wg.Done()
		for {
			var event Event

			select {
			case event = <-q.input:
			case <-q.doneCh:
				return
			}

			select {
			case q.output <- event:
			case <-q.doneCh:
				return
			}
		}
	}()
}

// Close closes the event queue.
func (q *Queue) Close() error {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.closed {
		return nil
	}

	close(q.doneCh)
	q.wg.Wait()

	close(q.output)

	finalErr := lb.RetryError{}

	if q.conn != nil {
		err := q.conn.Close()
		if err != nil {
			err = fmt.Errorf("failed to close connection: %w", err)
			finalErr.RawErrors = append(finalErr.RawErrors, err)
		}
	}

	q.closed = true

	if len(finalErr.RawErrors) > 0 {
		err := errors.New("unable to close queue")
		finalErr.RawErrors = append(finalErr.RawErrors, err)
		finalErr.Final = err
		return finalErr
	}

	return nil
}
