package config

import "os"

type Config struct {
	HTTPPort   string
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
}

func Load() Config {
	port := os.Getenv("HTTP_PORT")

	if port == "" {
		port = "8080"
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "127.0.0.1"
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "3307"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "orderpulse"
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "orderpulse"
	}

	dbPassword := os.Getenv("DB_PASSWORD")

	return Config{
		HTTPPort:   port,
		DBHost:     dbHost,
		DBPort:     dbPort,
		DBName:     dbName,
		DBUser:     dbUser,
		DBPassword: dbPassword,
	}
}
