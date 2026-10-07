package main

import (
	"log"
	"strconv"

	"github.com/Max-shiri-90/OrderPulse/internal/config"
	"github.com/Max-shiri-90/OrderPulse/internal/database"
	server "github.com/Max-shiri-90/OrderPulse/internal/http"
)

func main() {
	cfg := config.Load()

	port, err := strconv.Atoi(cfg.HTTPPort)
	if err != nil {
		log.Fatal("invalid HTTP_PORT:", err)
	}

	db, err := database.NewMySQL(database.Config{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		Name:     cfg.DBName,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
	})
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	defer db.Close()

	app := server.NewServer(port, db)

	log.Println("OrderPulse is running on :" + cfg.HTTPPort)

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
