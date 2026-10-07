package memory

import (
	"context"
	"sync"

	"github.com/Max-shiri-90/OrderPulse/internal/product"
)

type Repository struct {
	mu       sync.RWMutex
	products []product.Product
	nextID   int64
}

func NewRepository() *Repository {
	return &Repository{
		products: make([]product.Product, 0),
		nextID:   1,
	}
}

func (r *Repository) Create(ctx context.Context, p product.Product) (product.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	p.ID = r.nextID
	r.nextID++

	r.products = append(r.products, p)

	return p, nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (product.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, p := range r.products {
		if p.ID == id {
			return p, nil
		}
	}

	return product.Product{}, product.ErrNotFound
}

func (r *Repository) List(ctx context.Context) ([]product.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	products := make([]product.Product, len(r.products))
	copy(products, r.products)

	return products, nil
}
