package server

import (
	"fmt"
	"net/http"

	"github.com/Max-shiri-90/OrderPulse/internal/http/handler"
	"github.com/Max-shiri-90/OrderPulse/internal/product"
	"github.com/Max-shiri-90/OrderPulse/internal/product/memory"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(port int) *Server {
	mux := http.NewServeMux()

	healthHandler := handler.Health

	productRepository := memory.NewRepository()
	productService := product.NewService(productRepository)
	productHandler := handler.NewProductHandler(productService)

	mux.HandleFunc("/health", healthHandler)

	mux.HandleFunc("POST /api/v1/products", productHandler.Create)
	mux.HandleFunc("GET /api/v1/products", productHandler.List)
	mux.HandleFunc("GET /api/v1/products/{id}", productHandler.GetByID)

	return &Server{
		httpServer: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: mux,
		},
	}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}
