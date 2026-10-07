package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/Max-shiri-90/OrderPulse/internal/product"
)

type ProductHandler struct {
	service *product.Service
}

func NewProductHandler(service *product.Service) *ProductHandler {
	return &ProductHandler{
		service: service,
	}
}

type createProductRequest struct {
	Name     string `json:"name"`
	Price    int64  `json:"price"`
	Quantity int    `json:"quantity"`
}

type productResponse struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Price    int64  `json:"price"`
	Quantity int    `json:"quantity"`
}

func newProductResponse(p product.Product) productResponse {
	return productResponse{
		ID:       p.ID,
		Name:     p.Name,
		Price:    p.Price,
		Quantity: p.Quantity,
	}
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createProductRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	p, err := h.service.Create(
		r.Context(),
		req.Name,
		req.Price,
		req.Quantity,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := newProductResponse(p)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		return
	}
}

func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.List(r.Context())
	if err != nil {
		http.Error(w, "failed to list products", http.StatusInternalServerError)
		return
	}

	responses := make([]productResponse, 0, len(products))

	for _, p := range products {
		responses = append(responses, newProductResponse(p))
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(responses); err != nil {
		return
	}
}

func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	p, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, product.ErrNotFound) {
			http.Error(w, "product not found", http.StatusNotFound)
			return
		}

		http.Error(w, "failed to get product", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	response := newProductResponse(p)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		return
	}
}
