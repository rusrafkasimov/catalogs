package repository

import (
	"context"
	"fmt"
	"github.com/afiskon/promtail-client/promtail"
	"github.com/opentracing/opentracing-go"
	"github.com/rusrafkasimov/catalogs/internal/queue"
	"github.com/rusrafkasimov/catalogs/internal/trace"
	"github.com/rusrafkasimov/catalogs/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	companyCollection = "catalogs"
	mgoDatabase       = "dev"
)

type CatalogsRepository interface {
	CreateCatalog(ctx context.Context, model *models.Catalog, span opentracing.Span) (*models.Catalog, error)
	FindCatalogByID(ctx context.Context, id primitive.ObjectID, span opentracing.Span) (*models.Catalog, error)
	FindCatalogsByCategory(ctx context.Context, category string, span opentracing.Span) ([]*models.Catalog, error)
	FindCatalogsCategories(ctx context.Context, span opentracing.Span) ([]string, error)
	UpdateCatalog(ctx context.Context, id string, model *models.Catalog, span opentracing.Span) (*models.Catalog, error)
	DeleteCatalog(ctx context.Context, id string, span opentracing.Span) bool
}

func NewCatalogsRepository(ct *mongo.Client, eventQueue queue.EventQueue, logger promtail.Client) *CatalogsRepo {
	collection := ct.Database(mgoDatabase).Collection(companyCollection)
	return &CatalogsRepo{ct, collection, logger, eventQueue}
}

type CatalogsRepo struct {
	conn       *mongo.Client
	collection *mongo.Collection
	logger     promtail.Client
	eventQueue queue.EventQueue
}

func (m *CatalogsRepo) CreateCatalog(ctx context.Context, model *models.Catalog, span opentracing.Span) (*models.Catalog, error) {
	tracer := opentracing.GlobalTracer()
	repoSpan := tracer.StartSpan("Repo:CreateCatalog", opentracing.ChildOf(span.Context()))
	defer repoSpan.Finish()

	model.ID = primitive.NewObjectID()
	_, err := m.collection.InsertOne(ctx, model)
	if err != nil {
		return nil, err
	}

	op := &models.Operation{
		Type:    models.OperationTypeCatalogs,
		Method:  models.OperationMethodUpsert,
		Catalog: model,
	}

	if err = m.eventQueue.Publish(op); err != nil {
		return nil, err
	}

	return model, nil
}

func (m *CatalogsRepo) FindCatalogByID(ctx context.Context, id primitive.ObjectID, span opentracing.Span) (*models.Catalog, error) {
	tracer := opentracing.GlobalTracer()
	repoSpan := tracer.StartSpan("Repo:FindCatalogByID", opentracing.ChildOf(span.Context()))
	defer repoSpan.Finish()

	var newDocument *models.Catalog
	err := m.collection.FindOne(ctx, bson.D{{"_id", id}}).Decode(&newDocument)
	if err != nil {
		return nil, err
	}

	return newDocument, nil
}

func (m *CatalogsRepo) FindCatalogsByCategory(ctx context.Context, category string, span opentracing.Span) ([]*models.Catalog, error) {
	tracer := opentracing.GlobalTracer()
	repoSpan := tracer.StartSpan("Repo:FindCatalogsByCategory", opentracing.ChildOf(span.Context()))
	defer repoSpan.Finish()

	documents, err := m.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	var newDocuments []*models.Catalog
	err = documents.All(ctx, &newDocuments)

	return newDocuments, nil
}

func (m *CatalogsRepo) FindCatalogsCategories(ctx context.Context, span opentracing.Span) ([]string, error) {
	tracer := opentracing.GlobalTracer()
	repoSpan := tracer.StartSpan("Repo:FindCatalogsCategories", opentracing.ChildOf(span.Context()))
	defer repoSpan.Finish()

	var categories []string

	groupStage := bson.D{{"$group", bson.D{{"_id", "$category"}}}}

	docs, err := m.collection.Aggregate(ctx, mongo.Pipeline{groupStage})
	if err != nil {
		return categories, err
	}

	var categoriesRows []bson.M
	if err = docs.All(ctx, &categoriesRows); err != nil {
		panic(err)
	}

	for _, v := range categoriesRows {
		value, ok := v["_id"]
		if ok {
			categories = append(categories, value.(string))

		}
	}

	return categories, nil
}

func (m *CatalogsRepo) UpdateCatalog(ctx context.Context, id string, model *models.Catalog, span opentracing.Span) (*models.Catalog, error) {
	tracer := opentracing.GlobalTracer()
	repoSpan := tracer.StartSpan("Repo:UpdateCatalogs", opentracing.ChildOf(span.Context()))
	defer repoSpan.Finish()

	updatedId, _ := primitive.ObjectIDFromHex(id)
	model.ID = updatedId
	res, err := m.collection.ReplaceOne(ctx, bson.M{"_id": updatedId}, model)
	if err != nil {
		return nil, fmt.Errorf("failed to replace one: %w", err)
	}

	if res.MatchedCount == 0 {
		return nil, fmt.Errorf("not found record replace one: %w", err)
	}

	op := &models.Operation{
		Type:    models.OperationTypeCatalogs,
		Method:  models.OperationMethodUpsert,
		Catalog: model,
	}
	if err = m.eventQueue.Publish(op); err != nil {
		return nil, err
	}

	return model, nil
}

func (m *CatalogsRepo) DeleteCatalog(ctx context.Context, id string, span opentracing.Span) bool {
	tracer := opentracing.GlobalTracer()
	repoSpan := tracer.StartSpan("Repo:DeleteCatalogs", opentracing.ChildOf(span.Context()))
	defer repoSpan.Finish()

	var newDocument *models.Catalog

	deletedId, _ := primitive.ObjectIDFromHex(id)

	update := bson.M{
		"$set": bson.M{"active": false},
	}

	result := m.collection.FindOneAndUpdate(ctx, bson.D{{"_id", deletedId}}, update)
	if result.Err() != nil {
		trace.OnError(m.logger, span, result.Err())
		return false
	}

	err := result.Decode(&newDocument)
	if err != nil {
		trace.OnError(m.logger, span, err)
		return false
	}

	op := &models.Operation{
		Type:   models.OperationTypeCatalogs,
		Method: models.OperationMethodDelete,
		Catalog: &models.Catalog{
			ID: newDocument.ID,
		},
	}

	if err := m.eventQueue.Publish(op); err != nil {
		trace.OnError(m.logger, span, err)
		return false
	}

	return true
}
