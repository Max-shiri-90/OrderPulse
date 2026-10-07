package product

import (
	"context"
	"errors"
)

var (
	ErrInvalidName     = errors.New("product name cannot be empty")
	ErrInvalidPrice    = errors.New("product price must be greater than zero")
	ErrInvalidQuantity = errors.New("product quantity cannot be negative")
	ErrNotFound        = errors.New("product not found")
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Create(
	ctx context.Context,
	name string,
	price int64,
	quantity int,
) (Product, error) {
	if name == "" {
		return Product{}, ErrInvalidName
	}

	if price <= 0 {
		return Product{}, ErrInvalidPrice
	}

	if quantity < 0 {
		return Product{}, ErrInvalidQuantity
	}

	p := Product{
		Name:     name,
		Price:    price,
		Quantity: quantity,
	}

	return s.repository.Create(ctx, p)
}

func (s *Service) GetByID(
	ctx context.Context,
	id int64,
) (Product, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *Service) List(
	ctx context.Context,
) ([]Product, error) {
	return s.repository.List(ctx)
}
