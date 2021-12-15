package memstore

import (
	"context"
	"github.com/rusrafkasimov/catalogs/pkg/models"
	"sync"
)

type MemStore interface {
	UpsertCatalog(catalog *models.Catalog)
	UpsertCatalogByCategory (catalog *models.Catalog)
	GetCatalog(id string) (*models.Catalog, bool)
	GetCatalogs(ids []string, sorted bool) []*models.Catalog
	GetCategories() []string
	GetRefsByCategoryAndQuery(category string, query string, sorted bool) ([]*models.Catalog, bool)
	RemoveCatalog(id string)
}

type memStore struct {
	context context.Context
	catalog struct {
		sync.RWMutex
		data map[string]*models.Catalog
		category map[string]map[string]bool
	}
}

func NewMemStore(context context.Context) MemStore {
	m := &memStore{
		context: context,
	}
	m.catalog.data = make(map[string]*models.Catalog)
	m.catalog.category = make(map[string]map[string]bool)
	return m
}

func (m *memStore) UpsertCatalog(catalog *models.Catalog) {

}

func (m *memStore) UpsertCatalogByCategory (catalog *models.Catalog) {

}

func (m *memStore) GetCatalog(id string) (*models.Catalog, bool) {

}

func (m *memStore) GetCatalogs(ids []string, sorted bool) []*models.Catalog {

}

func (m *memStore) GetCategories() []string {

}

func (m *memStore) GetCatalogByCategoryAndQuery(category string, query string, sorted bool) ([]*models.Catalog, bool) {

}

func (m *memStore) RemoveCatalog(id string) {

}

