package product

import "context"

type Repository interface {
	Create(ctx context.Context, p Product) (Product, error)
	GetByID(ctx context.Context, id int64) (Product, error)
	List(ctx context.Context) ([]Product, error)
}
