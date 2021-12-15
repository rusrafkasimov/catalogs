package repository

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/rusrafkasimov/catalogs/internal/queue"
	"github.com/rusrafkasimov/catalogs/pkg/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	companyCollection  = "catalogs"
	mgoDatabase = "dev"
)

type CatalogsRepository interface {
	CreateCatalog(ctx context.Context, model *models.Catalog, span opentracing.Span) (*models.Catalog, error)
	FindCatalogByID(ctx context.Context, id primitive.ObjectID, span opentracing.Span) (*models.Catalog, error)
	FindCatalogsByCategory(ctx context.Context, category string, span opentracing.Span) ([]*models.Catalog, error)
	FindCatalogsCategories(ctx context.Context, span opentracing.Span) ([]string, error)
	UpdateCatalogs(ctx context.Context, id string, model *models.Catalog, span opentracing.Span) (*models.Catalog, error)
	DeleteCatalogs(ctx context.Context, id string, span opentracing.Span) bool
}

func NewCatalogsRepository(ct *mongo.Client, eventQueue queue.EventQueue) *CatalogsRepo {
	collection := ct.Database(mgoDatabase).Collection(companyCollection)
	return &CatalogsRepo{ct, collection, eventQueue}
}

type CatalogsRepo struct {
	conn       *mongo.Client
	collection *mongo.Collection
	eventQueue queue.EventQueue
}

func (m *CatalogsRepo) CreateCatalog(ctx context.Context, model *models.Catalog, span opentracing.Span) (*models.Catalog, error) {

}

func (m *CatalogsRepo) FindCatalogByID(ctx context.Context, id primitive.ObjectID, span opentracing.Span) (*models.Catalog, error) {

}

func (m *CatalogsRepo) FindCatalogsByCategory(ctx context.Context, category string, span opentracing.Span) ([]*models.Catalog, error) {

}

func (m *CatalogsRepo) FindCatalogsCategories(ctx context.Context, span opentracing.Span) ([]string, error) {

}

func (m *CatalogsRepo) UpdateCatalogs(ctx context.Context, id string, model *models.Catalog, span opentracing.Span) (*models.Catalog, error) {

}

func (m *CatalogsRepo) DeleteCatalogs(ctx context.Context, id string, span opentracing.Span) bool {

}


