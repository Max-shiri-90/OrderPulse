package product_test

import (
	"context"
	"testing"

	"github.com/Max-shiri-90/OrderPulse/internal/product"
	"github.com/Max-shiri-90/OrderPulse/internal/product/memory"
)

func TestCreateProduct(t *testing.T) {
	repository := memory.NewRepository()
	service := product.NewService(repository)

	p, err := service.Create(
		context.Background(),
		"Mechanical Keyboard",
		7500000,
		25,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if p.ID != 1 {
		t.Fatalf("expected ID 1, got %d", p.ID)
	}

	if p.Name != "Mechanical Keyboard" {
		t.Fatalf("expected product name %q, got %q", "Mechanical Keyboard", p.Name)
	}

	if p.Price != 7500000 {
		t.Fatalf("expected price 7500000, got %d", p.Price)
	}

	if p.Quantity != 25 {
		t.Fatalf("expected quantity 25, got %d", p.Quantity)
	}
}

func TestCreateProductValidation(t *testing.T) {
	repository := memory.NewRepository()
	service := product.NewService(repository)

	_, err := service.Create(
		context.Background(),
		"",
		7500000,
		25,
	)

	if err != product.ErrInvalidName {
		t.Fatalf("expected ErrInvalidName, got %v", err)
	}

	_, err = service.Create(
		context.Background(),
		"Keyboard",
		0,
		25,
	)

	if err != product.ErrInvalidPrice {
		t.Fatalf("expected ErrInvalidPrice, got %v", err)
	}

	_, err = service.Create(
		context.Background(),
		"Keyboard",
		7500000,
		-1,
	)

	if err != product.ErrInvalidQuantity {
		t.Fatalf("expected ErrInvalidQuantity, got %v", err)
	}
}
