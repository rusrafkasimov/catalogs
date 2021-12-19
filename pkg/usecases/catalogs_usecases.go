package usecases

import (
	"context"
	"errors"
	"github.com/afiskon/promtail-client/promtail"
	"github.com/mitchellh/mapstructure"
	"github.com/opentracing/opentracing-go"
	"github.com/rusrafkasimov/catalogs/internal/trace"
	"github.com/rusrafkasimov/catalogs/pkg/dto"
	"github.com/rusrafkasimov/catalogs/pkg/models"
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
	tracer := opentracing.GlobalTracer()
	useCaseSpan := tracer.StartSpan("UCase:CreateCatalog", opentracing.ChildOf(span.Context()))
	defer useCaseSpan.Finish()
	var result dto.CreateCatalogResponse
	var catalog *models.Catalog

	err := mapstructure.Decode(request, &catalog)
	if err != nil {
		trace.OnError(c.logger, useCaseSpan, err)
		return nil, err
	}

	model, err := c.rep.CreateCatalog(ctx, catalog, useCaseSpan)
	if err != nil {
		trace.OnError(c.logger, useCaseSpan, err)
		return nil, err
	}

	err = mapstructure.Decode(model, &result.Payload)
	if err != nil {
		trace.OnError(c.logger, useCaseSpan, err)
		return &result, err
	}
	return &result, nil
}

func (c *CatalogsUC) GetCatalogs(ctx context.Context, request *dto.CatalogsRequest, span opentracing.Span) (*dto.GetCatalogsResponse, error) {
	tracer := opentracing.GlobalTracer()
	useCaseSpan := tracer.StartSpan("UCase:GetCatalogs", opentracing.ChildOf(span.Context()))
	defer useCaseSpan.Finish()
	var result dto.GetCatalogsResponse

	documents, ok := c.store.GetCatalogByCategoryAndQuery(request.Category, request.Query, request.Sorted)
	if !ok {
		trace.OnError(c.logger, useCaseSpan, errors.New("memstore is empty"))
		return &result, errors.New("memstore is empty")
	}

	for _, entry := range documents {
		newModel := dto.CatalogResponse{
			ID:       entry.ID,
			Name:     entry.Name,
			Category: entry.Category,
			Value:    entry.Value,
		}
		result.Payload = append(result.Payload, newModel)
	}

	return &result, nil
}

func (c *CatalogsUC) GetCatalogCategories(ctx context.Context, span opentracing.Span) (*dto.GetCategoriesResponse, error) {
	tracer := opentracing.GlobalTracer()
	useCaseSpan := tracer.StartSpan("UCase:GetCatalogCategories", opentracing.ChildOf(span.Context()))
	defer useCaseSpan.Finish()
	var result dto.GetCategoriesResponse

	categories := c.store.GetCategories()

	result.Payload = categories

	return &result, nil
}

func (c *CatalogsUC) GetCatalogByID(ctx context.Context, id string, span opentracing.Span) (*dto.GetCatalogResponse, error) {
	tracer := opentracing.GlobalTracer()
	useCaseSpan := tracer.StartSpan("UCase:GetCatalogByID", opentracing.ChildOf(span.Context()))
	defer useCaseSpan.Finish()
	var result dto.GetCatalogResponse

	media, ok := c.store.GetCatalog(id)
	if !ok {
		trace.OnError(c.logger, useCaseSpan, errors.New("memstore error"))
		return nil, errors.New("not found")
	}

	result.Payload.ID = media.ID
	result.Payload.Active = media.Active
	result.Payload.Name = media.Name
	result.Payload.Value = media.Value
	result.Payload.Category = media.Category

	return &result, nil
}

func (c *CatalogsUC) UpdateCatalogByID(ctx context.Context, request *dto.CatalogRequest, span opentracing.Span) (*dto.UpdateCatalogResponse, error) {
	tracer := opentracing.GlobalTracer()
	useCaseSpan := tracer.StartSpan("UCase:UpdateCatalogByID", opentracing.ChildOf(span.Context()))
	defer useCaseSpan.Finish()
	var result dto.UpdateCatalogResponse
	var catalog *models.Catalog

	err := mapstructure.Decode(request, &catalog)
	if err != nil {
		trace.OnError(c.logger, useCaseSpan, err)
		return nil, err
	}

	model, err := c.rep.UpdateCatalog(ctx, request.ID, catalog, useCaseSpan)
	if err != nil {
		trace.OnError(c.logger, useCaseSpan, err)
		return nil, err
	}

	err = mapstructure.Decode(model, &result.Payload)
	if err != nil {
		trace.OnError(c.logger, useCaseSpan, err)
		return &result, err
	}

	return &result, nil
}

func (c *CatalogsUC) DeleteCatalogByID(ctx context.Context, id string, span opentracing.Span) (*dto.DeleteCatalogResponse, error) {
	tracer := opentracing.GlobalTracer()
	useCaseSpan := tracer.StartSpan("UCase:DeleteCatalogByID", opentracing.ChildOf(span.Context()))
	defer useCaseSpan.Finish()
	var result dto.DeleteCatalogResponse

	ok := c.rep.DeleteCatalog(ctx, id, useCaseSpan)
	if !ok {
		trace.OnError(c.logger, useCaseSpan, errors.New("not found"))
		return nil, errors.New("not found")
	}

	result.Payload.Active = false

	return &result, nil
}
