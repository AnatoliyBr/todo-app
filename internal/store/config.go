package store

import (
	"os"
)

type Config struct {
	DatabaseURL string
}

func NewConfig() *Config {
	databaseURL := os.Getenv("DATABASE_URL_LOCALHOST")

	databaseURL += "?sslmode=disable"

	if databaseURL == "" {
		databaseURL = "postgres://dev:qwerty@localhost:5432/todo_dev"
	}

	return &Config{
		DatabaseURL: databaseURL,
	}
}
