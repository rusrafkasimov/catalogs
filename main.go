package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/rusrafkasimov/catalogs/internal/config"
	"github.com/rusrafkasimov/catalogs/internal/logger"
	"github.com/rusrafkasimov/catalogs/internal/mongo"
	"github.com/rusrafkasimov/catalogs/internal/queue"
	"github.com/rusrafkasimov/catalogs/internal/trace"
	"github.com/rusrafkasimov/catalogs/internal/vault"
	"github.com/rusrafkasimov/catalogs/pkg/delivery/router"
	"time"
)

const (
	Name           = "Catalogs"
	contextKeyName = "Name"
	serverTimeout  = 10 * time.Second
)

func main() {
	ctx := context.Background()

	id, err := gonanoid.New()
	if err != nil {
		fmt.Println("Can't generate new node ID")
		return
	}

	ctx = context.WithValue(
		ctx,
		contextKeyName,
		Name+"_"+id,
	)

	var env string

	flag.StringVar(&env, "env", ".env.local", "Environment Variables filename")
	flag.Parse()

	// Load service configuration from environment
	if err := config.LoadConfig(env); err != nil {
		fmt.Printf("Error: can't load env. %" + err.Error())
	}

	// Initialize vault
	vaultProvider := vault.NewVaultProvider()
	appConfig := config.NewConfig(vaultProvider)

	// Initialize logger
	loki, err := logger.NewLogger(Name, "api", appConfig)
	if err != nil {
		fmt.Println("Error while connect to loki")
	}

	// Initialize tracing
	closer, err := trace.InitJaegerTracing(ctx, contextKeyName, appConfig)
	if err != nil {
		loki.Errorf("Error while init tracing")
	}
	defer closer.Close()

	// Initialize Database
	mgoDB, err := mongo.InitDatabase(ctx, loki, appConfig)
	if err != nil {
		loki.Errorf("Error init database")
	}

	// Initialize NATS Queue
	newQueue, err := queue.NewQueue(ctx, loki, appConfig)
	if err != nil {
		loki.Errorf("Error init new queue")
	}

	// Build context
	repoCtx := router.BuildRepositoryContext(mgoDB.Client, ctx, newQueue, loki)
	ucCtx := router.BuildUcaseContext(repoCtx, loki)
	appCtx := router.BuildApplicationContext(ucCtx, loki)

	// Initialize gin routes and run server
	rGin := gin.Default()
	gin.ForceConsoleColor()
	router.MapUrl(rGin, appCtx)

	httpHost, err := appConfig.Get("HTTP_HOST")
	if err != nil {
		loki.Errorf(err.Error())
	}

	httpPort := ":8090"

	if err = rGin.Run(httpPort); err != nil {
		loki.Errorf("Error: can't start GIN router. %s", err.Error())
	}

	loki.Infof("Upstream started at %v", httpHost+httpPort)

	defer func(){
		loki.Infof("Catalogs service stopped")
		_ = newQueue.Close()
	}()
}
