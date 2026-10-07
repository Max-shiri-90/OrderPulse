package main

import (
	"log"

	server "github.com/Max-shiri-90/OrderPulse/internal/http"
)

func main() {
	app := server.NewServer(8080)

	log.Println("OrderPulse is running on :8080")

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}