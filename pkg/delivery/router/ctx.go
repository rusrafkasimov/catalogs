package router

import (
	"context"
	"github.com/afiskon/promtail-client/promtail"
	"github.com/rusrafkasimov/catalogs/internal/queue"
	"github.com/rusrafkasimov/catalogs/pkg/controllers"
	"github.com/rusrafkasimov/catalogs/pkg/repository/memstore"
	repository "github.com/rusrafkasimov/catalogs/pkg/repository/mongo"
	"github.com/rusrafkasimov/catalogs/pkg/usecases"
	"go.mongodb.org/mongo-driver/mongo"
)

type RepositoryContext struct {
	CatalogRep *repository.CatalogsRepo
	CatalogMem memstore.MemStore
}

type UseCaseContext struct {
	catUseCases *usecases.CatalogsUC
}

type ApplicationContext struct {
	CatalogsController *controllers.CatalogsController
}

func BuildRepositoryContext(mgo *mongo.Client, ctx context.Context, eq queue.EventQueue, logger promtail.Client) *RepositoryContext {
	return &RepositoryContext{
		CatalogRep: repository.NewCatalogsRepository(mgo, eq, logger),
		CatalogMem: memstore.NewMemStore(ctx),
	}
}

func BuildUcaseContext(repoCtx *RepositoryContext, logger promtail.Client) *UseCaseContext {
	return &UseCaseContext{
		catUseCases: usecases.NewCatalogsUseCases(repoCtx.CatalogRep, repoCtx.CatalogMem, logger),
	}
}

func BuildApplicationContext(ucCtx *UseCaseContext, logger promtail.Client) *ApplicationContext {
	return &ApplicationContext{
		CatalogsController: controllers.NewCatalogsController(ucCtx.catUseCases, logger),
	}
}



