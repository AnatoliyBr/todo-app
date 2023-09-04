package store

import (
	"os"
)

type Config struct {
	DatabaseURL string
}

func NewConfig() *Config {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://dev:qwerty@localhost:5432/todo_dev"
	}

	databaseURL += "?sslmode=disable"

	return &Config{
		DatabaseURL: databaseURL,
	}
}
