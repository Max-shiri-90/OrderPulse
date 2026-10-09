package server

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/Max-shiri-90/OrderPulse/internal/http/handler"
	"github.com/Max-shiri-90/OrderPulse/internal/http/middleware"
	"github.com/Max-shiri-90/OrderPulse/internal/product"
	productmysql "github.com/Max-shiri-90/OrderPulse/internal/product/mysql"
	"github.com/Max-shiri-90/OrderPulse/internal/user"
	usermysql "github.com/Max-shiri-90/OrderPulse/internal/user/mysql"
)

type Server struct {
	httpServer *http.Server
	jwtSecret  string
}

func NewServer(port int, db *sql.DB, jwtSecret string) *Server {
	mux := http.NewServeMux()

	// Health
	healthHandler := handler.Health

	// Products
	productRepository := productmysql.NewRepository(db)
	productService := product.NewService(productRepository)
	productHandler := handler.NewProductHandler(productService)

	// Users / Authentication
	userRepository := usermysql.NewRepository(db)
	userService := user.NewService(userRepository)
	authHandler := handler.NewAuthHandler(userService, jwtSecret)

	// Routes
	mux.HandleFunc("/health", healthHandler)

	mux.HandleFunc(
		"POST /api/v1/products",
		productHandler.Create,
	)

	mux.HandleFunc(
		"GET /api/v1/products",
		productHandler.List,
	)

	mux.HandleFunc(
		"GET /api/v1/products/{id}",
		productHandler.GetByID,
	)

	mux.HandleFunc(
		"POST /api/v1/auth/register",
		authHandler.Register,
	)

	mux.HandleFunc(
		"POST /api/v1/auth/login",
		authHandler.Login,
	)

	mux.Handle(
		"GET /api/v1/me",
		middleware.RequireAuth(
			jwtSecret,
			http.HandlerFunc(authHandler.GetMe),
		),
	)

	return &Server{
		httpServer: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: mux,
		},
		jwtSecret: jwtSecret,
	}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}
