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
		databaseURL = "host=localhost user=dev password=qwerty dbname=todo_test sslmode=disable"
	}

	return &Config{
		DatabaseURL: databaseURL,
	}
}
