package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Max-shiri-90/OrderPulse/internal/product"
	"github.com/Max-shiri-90/OrderPulse/internal/product/memory"
)

func newTestProductHandler() *ProductHandler {
	repository := memory.NewRepository()
	service := product.NewService(repository)

	return NewProductHandler(service)
}

func TestProductHandlerCreate(t *testing.T) {
	handler := newTestProductHandler()

	body := `{"name":"Mechanical Keyboard","price":7500000,"quantity":25}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/products",
		strings.NewReader(body),
	)

	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	expected := `{"id":1,"name":"Mechanical Keyboard","price":7500000,"quantity":25}` + "\n"

	if rec.Body.String() != expected {
		t.Fatalf("expected body %q, got %q", expected, rec.Body.String())
	}
}

func TestProductHandlerCreateInvalidName(t *testing.T) {
	handler := newTestProductHandler()

	body := `{"name":"","price":7500000,"quantity":25}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/products",
		strings.NewReader(body),
	)

	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestProductHandlerGetByID(t *testing.T) {
	repository := memory.NewRepository()
	service := product.NewService(repository)
	handler := NewProductHandler(service)

	_, err := service.Create(
		context.Background(),
		"Mechanical Keyboard",
		7500000,
		25,
	)

	if err != nil {
		t.Fatalf("failed to create test product: %v", err)
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/products/1",
		nil,
	)

	req.SetPathValue("id", "1")

	rec := httptest.NewRecorder()

	handler.GetByID(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	expected := `{"id":1,"name":"Mechanical Keyboard","price":7500000,"quantity":25}` + "\n"

	if rec.Body.String() != expected {
		t.Fatalf("expected body %q, got %q", expected, rec.Body.String())
	}
}

func TestProductHandlerGetByIDNotFound(t *testing.T) {
	handler := newTestProductHandler()

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/products/999",
		nil,
	)

	req.SetPathValue("id", "999")

	rec := httptest.NewRecorder()

	handler.GetByID(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}
