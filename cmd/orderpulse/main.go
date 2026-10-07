package main

import (
	"log"
	"strconv"

	"github.com/Max-shiri-90/OrderPulse/internal/config"
	server "github.com/Max-shiri-90/OrderPulse/internal/http"
)

func main() {
	cfg := config.Load()

	port, err := strconv.Atoi(cfg.HTTPPort)
	if err != nil {
		log.Fatal("invalid HTTP_PORT:", err)
	}

	app := server.NewServer(port)

	log.Println("OrderPulse is running on :" + cfg.HTTPPort)

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
