package replicator

import (
	"context"
	"errors"
	"fmt"
	"github.com/rusrafkasimov/catalogs/internal/queue"
	"github.com/rusrafkasimov/catalogs/internal/trace"
	"github.com/rusrafkasimov/catalogs/pkg/models"
	"github.com/rusrafkasimov/catalogs/pkg/repository/memstore"
	repository "github.com/rusrafkasimov/catalogs/pkg/repository/mongo"
	"sync/atomic"

	"github.com/afiskon/promtail-client/promtail"
	"github.com/opentracing/opentracing-go"
)

const (
	OperationTypeCatalogs models.OperationType = "catalogs"
	OperationMethodUpsert models.OperationMethod = "upsert"
	OperationMethodDelete models.OperationMethod = "delete"
)

type Replicator interface {
	Replicate() error
	loadDataFromStorage(ctx context.Context) error
	loadCatalogs(ctx context.Context) error
}

type replicator struct {
	ctx            context.Context
	repo           repository.CatalogsRepository
	memStore       memstore.MemStore
	eventQueue     queue.EventQueue
	logger         promtail.Client
	ready          uint32
	operationTypes map[models.OperationType]bool
}

func New(ctx context.Context, repo repository.CatalogsRepository, memStore memstore.MemStore, eventQueue queue.EventQueue, logger promtail.Client, operationTypes []models.OperationType) *replicator {

	types := make(map[models.OperationType]bool, len(operationTypes))
	for _, ot := range operationTypes {
		types[ot] = true
	}

	return &replicator{ctx: ctx, repo: repo, memStore: memStore, eventQueue: eventQueue, logger: logger, operationTypes: types}
}

// setReady set atomic ready value
func (r *replicator) setReady(val bool) {
	if val {
		atomic.StoreUint32(&r.ready, 1)
	} else {
		atomic.StoreUint32(&r.ready, 0)
	}
}

// Ready returns true if the replicator successfully completed reading catalog data from
// the storage and started to handle replication events from the event queue.
func (r *replicator) Ready() bool {
	return atomic.LoadUint32(&r.ready) == 1
}

// Replicate start loading data from storage and check replication events
func (r *replicator) Replicate(ctx context.Context) error {
	err := r.loadDataFromStorage(ctx)
	r.setReady(true)

	if err := r.handleReplicationEvents(ctx); err != nil {
		return fmt.Errorf("failed to handle replication events from event queue: %w", err)
	}

	return err
}

// loadDataFromStorage start loading all data to storage
func (r *replicator) loadDataFromStorage(ctx context.Context) error {
	if err := r.loadCatalogs(ctx); err != nil {
		return fmt.Errorf("failed to load data: %w", err)
	}

	return nil
}

// loadCatalogs load catalogs from DB, send all data to memory storage
func (r *replicator) loadCatalogs(ctx context.Context) error {
	if !r.operationTypes[OperationTypeCatalogs] {
		return nil
	}

	tracer := opentracing.GlobalTracer()
	replicatorSpan := tracer.StartSpan("Replicator:LoadCatalogs")
	defer replicatorSpan.Finish()

	refBooks, err := r.repo.FindCatalogsByCategory(ctx, "", replicatorSpan)
	replicatorSpan.SetTag("Count catalogs for replicate", len(refBooks))
	if err != nil {
		trace.OnError(r.logger, replicatorSpan, err)
		return err
	}

	for _, entry := range refBooks {
		err = r.processOperation(&models.Operation{
			Type:    OperationTypeCatalogs,
			Method:  OperationMethodUpsert,
			Catalog: entry,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// processOperation depending on the operation type make changes to the memory storage
func (r *replicator) processOperation(op *models.Operation) (err error) {

	if op.Catalog == nil {
		return nil
	}

	switch op.Type {
	case OperationTypeCatalogs:
		switch op.Method {
		case OperationMethodDelete:
			r.memStore.RemoveCatalog(op.Catalog.ID.String())

		case OperationMethodUpsert:
			r.memStore.UpsertCatalog(op.Catalog)
			r.memStore.UpsertCatalogByCategory(op.Catalog)
		}
	}

	return nil
}

// handleReplicationEvents subscribe service on replication events and handle events
func (r *replicator) handleReplicationEvents(ctx context.Context) error {
	eventCh, err := r.eventQueue.Subscribe()
	if err != nil {
		return fmt.Errorf("failed to subscribe event queue: %w", err)
	}

	for {
		select {
		case evt, ok := <-eventCh:
			if !ok {
				return errors.New("event channel is closed")
			}

			err := r.processOperation(evt.Operation())
			if err != nil {
				trace.OnError(r.logger, nil, err)
				continue
			}

			err = evt.Ack()
			if err != nil {
				trace.OnError(r.logger, nil, err)
				continue
			}

		case <-ctx.Done():
			return nil
		}
	}
}
