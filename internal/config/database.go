package config

import (
	"fmt"
	"os"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string
}

func LoadDatabaseConfig() (*DatabaseConfig, error) {
	appEnv := getEnv("APP_ENV", "local")

	cfg := &DatabaseConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Name:     os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	// Allow very limited defaults ONLY for local/dev convenience
	if appEnv == "local" {
		if cfg.Host == "" {
			cfg.Host = "localhost"
		}
		if cfg.Port == "" {
			cfg.Port = "5432"
		}
	}


	return cfg, nil
}

func (c *DatabaseConfig) Redacted() string {
	return fmt.Sprintf(
		"postgres://%s:***@%s:%s/%s?sslmode=%s",
		c.User,
		c.Host,
		c.Port,
		c.Name,
		c.SSLMode,
	)
}

func (c *DatabaseConfig) buildDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.Name,
		c.SSLMode,
	)
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
