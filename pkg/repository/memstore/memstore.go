package memstore

import (
	"context"
	"github.com/rusrafkasimov/catalogs/pkg/models"
	"sort"
	"strings"
	"sync"
)

type MemStore interface {
	UpsertCatalog(catalog *models.Catalog)
	UpsertCatalogByCategory(catalog *models.Catalog)
	GetCatalog(id string) (*models.Catalog, bool)
	GetCatalogs(ids []string, sorted bool) []*models.Catalog
	GetCategories() []string
	GetCatalogByCategoryAndQuery(category string, query string, sorted bool) ([]*models.Catalog, bool)
	RemoveCatalog(id string)
}

type memStore struct {
	context context.Context
	catalog struct {
		sync.RWMutex
		data     map[string]*models.Catalog
		category map[string]map[string]bool
	}
}

func NewMemStore(context context.Context) *memStore {
	m := &memStore{
		context: context,
	}
	m.catalog.data = make(map[string]*models.Catalog)
	m.catalog.category = make(map[string]map[string]bool)
	return m
}

func (m *memStore) UpsertCatalog(catalog *models.Catalog) {
	m.catalog.Lock()
	defer m.catalog.Unlock()

	m.catalog.data[catalog.ID.String()] = catalog
}

func (m *memStore) UpsertCatalogByCategory(catalog *models.Catalog) {
	m.catalog.Lock()
	defer m.catalog.Unlock()

	item, ok := m.catalog.category[catalog.Category]
	if !ok {
		m.catalog.category[catalog.Category] = make(map[string]bool)
	}

	if item[catalog.ID.String()] {
		return
	}

	m.catalog.category[catalog.Category][catalog.ID.String()] = catalog.Active
}

func (m *memStore) GetCatalog(id string) (*models.Catalog, bool) {
	m.catalog.RLock()
	defer m.catalog.RUnlock()

	p, ok := m.catalog.data[id]

	return p, ok
}

func (m *memStore) GetCatalogs(ids []string, sorted bool) []*models.Catalog {
	m.catalog.RLock()
	defer m.catalog.RUnlock()

	out := make([]*models.Catalog, 0, len(ids))
	for _, id := range ids {
		p, ok := m.catalog.data[id]
		if ok {
			out = append(out, p)
		}
	}

	if sorted {
		sort.Sort(models.ByName(out))
	}

	return out
}

func (m *memStore) GetCategories() []string {
	m.catalog.RLock()
	defer m.catalog.RUnlock()

	var out []string

	for category, _ := range m.catalog.category {
		out = append(out, category)
	}

	return out
}

func (m *memStore) GetCatalogByCategoryAndQuery(category string, query string, sorted bool) ([]*models.Catalog, bool) {
	m.catalog.RLock()
	defer m.catalog.RUnlock()

	findedCategory, ok := m.catalog.category[category]
	var ids []string

	for key, _ := range findedCategory {
		ids = append(ids, key)
	}

	refs := m.GetCatalogs(ids, sorted)
	var filtered []*models.Catalog

	if query != "" {
		for _, ref := range refs {
			if strings.Contains(strings.ToLower(ref.Name), strings.ToLower(query)) {
				filtered = append(filtered, ref)
			}
		}
		return filtered, ok
	}

	return refs, ok
}

func (m *memStore) RemoveCatalog(id string) {
	m.catalog.Lock()
	defer m.catalog.Unlock()

	findedItem, ok := m.catalog.data[id]
	if !ok {
		return
	}

	delete(m.catalog.data, id)
	delete(m.catalog.category[findedItem.Category], id)
}
