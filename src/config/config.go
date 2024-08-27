package config

import (
	"errors"
	"os"
)

type Config struct {
	DatabaseURL string
}

func New() (*Config, error) {
	dbURL, found := os.LookupEnv("DATABASE_URL")
	if !found {
		return nil, errors.New("DATABASE_URL not found")
	}

	return &Config{
		DatabaseURL: dbURL,
	}, nil
}
