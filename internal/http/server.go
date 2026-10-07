package server

import (
	"fmt"
	"net/http"

	"github.com/Max-shiri-90/OrderPulse/internal/http/handler"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(port int) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", handler.Health)

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
