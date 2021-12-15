package usecases

import (
	"context"
	"github.com/afiskon/promtail-client/promtail"
	"github.com/opentracing/opentracing-go"
	"github.com/rusrafkasimov/catalogs/pkg/dto"
	"github.com/rusrafkasimov/catalogs/pkg/repository/memstore"
	repository "github.com/rusrafkasimov/catalogs/pkg/repository/mongo"
)

type CatalogsUseCase interface {
	CreateCatalog(ctx context.Context, request *dto.CatalogRequest, span opentracing.Span) (*dto.CreateCatalogResponse, error)
	GetCatalogs(ctx context.Context, request *dto.CatalogsRequest, span opentracing.Span) (*dto.GetCatalogsResponse, error)
	GetCatalogCategories(ctx context.Context, span opentracing.Span) (*dto.GetCategoriesResponse, error)
	GetCatalogByID(ctx context.Context, id string, span opentracing.Span) (*dto.GetCatalogResponse, error)
	UpdateCatalogByID(ctx context.Context, request *dto.CatalogRequest, span opentracing.Span) (*dto.UpdateCatalogResponse, error)
	DeleteCatalogByID(ctx context.Context, id string, span opentracing.Span) (*dto.DeleteCatalogResponse, error)
}

type CatalogsUC struct {
	rep    repository.CatalogsRepository
	store  memstore.MemStore
	logger promtail.Client
}

func NewCatalogsUseCases(rep repository.CatalogsRepository, store memstore.MemStore, logger promtail.Client) *CatalogsUC {
	return &CatalogsUC{
		rep:   rep,
		store: store,
		logger: logger,
	}
}

func (c *CatalogsUC) CreateCatalog(ctx context.Context, request *dto.CatalogRequest, span opentracing.Span) (*dto.CreateCatalogResponse, error) {

}

func (c *CatalogsUC) GetCatalogs(ctx context.Context, request *dto.CatalogsRequest, span opentracing.Span) (*dto.GetCatalogsResponse, error) {

}

func (c *CatalogsUC) GetCatalogCategories(ctx context.Context, span opentracing.Span) (*dto.GetCategoriesResponse, error) {

}

func (c *CatalogsUC) GetCatalogByID(ctx context.Context, id string, span opentracing.Span) (*dto.GetCatalogResponse, error) {

}

func (c *CatalogsUC) UpdateCatalogByID(ctx context.Context, request *dto.CatalogRequest, span opentracing.Span) (*dto.UpdateCatalogResponse, error) {

}

func (c *CatalogsUC) DeleteCatalogByID(ctx context.Context, id string, span opentracing.Span) (*dto.DeleteCatalogResponse, error) {

}
