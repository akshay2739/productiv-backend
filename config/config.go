package config

import (
	"fmt"
	"os"
)

// Config holds application configuration.
type Config struct {
	Port        string
	DatabaseURL string
	CORSOrigin  string
}

// Load reads configuration from environment variables with sensible defaults.
func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/discipline?sslmode=disable"),
		CORSOrigin:  getEnv("CORS_ORIGIN", "http://localhost:3000"),
	}
}

// DBConnString returns the formatted database connection string.
func (c *Config) DBConnString() string {
	return c.DatabaseURL
}

// ServerAddr returns the formatted server address.
func (c *Config) ServerAddr() string {
	return fmt.Sprintf(":%s", c.Port)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
