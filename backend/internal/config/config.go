package config

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	AppEnv   string
	HTTPAddr string
	Database Database
}

type Database struct {
	Host         string
	Port         uint16
	Name         string
	User         string
	PasswordFile string
	SSLMode      string
}

func Load() (Config, error) {
	appEnv := valueOrDefault("APP_ENV", "development")
	if appEnv != "development" && appEnv != "test" && appEnv != "production" {
		return Config{}, errors.New("APP_ENV must be development, test, or production")
	}

	portValue := valueOrDefault("DATABASE_PORT", "5432")
	port, err := strconv.ParseUint(portValue, 10, 16)
	if err != nil || port == 0 {
		return Config{}, errors.New("DATABASE_PORT must be an integer between 1 and 65535")
	}

	sslMode, exists := os.LookupEnv("DATABASE_SSLMODE")
	if !exists {
		if appEnv != "development" {
			return Config{}, errors.New("DATABASE_SSLMODE must be explicitly set outside development")
		}
		sslMode = "disable"
	}
	if appEnv == "production" && sslMode != "verify-full" {
		return Config{}, errors.New("DATABASE_SSLMODE=verify-full is required in production")
	}
	switch sslMode {
	case "disable", "allow", "prefer", "require", "verify-ca", "verify-full":
	default:
		return Config{}, errors.New("DATABASE_SSLMODE is not a supported PostgreSQL SSL mode")
	}

	cfg := Config{
		AppEnv:   appEnv,
		HTTPAddr: valueOrDefault("HTTP_ADDR", ":8080"),
		Database: Database{
			Host:         valueOrDefault("DATABASE_HOST", "127.0.0.1"),
			Port:         uint16(port),
			Name:         required("DATABASE_NAME"),
			User:         required("DATABASE_USER"),
			PasswordFile: required("DATABASE_PASSWORD_FILE"),
			SSLMode:      sslMode,
		},
	}

	if cfg.Database.Name == "" || cfg.Database.User == "" || cfg.Database.PasswordFile == "" {
		return Config{}, errors.New("DATABASE_NAME, DATABASE_USER, and DATABASE_PASSWORD_FILE are required")
	}
	if _, _, err := net.SplitHostPort(cfg.HTTPAddr); err != nil {
		return Config{}, fmt.Errorf("HTTP_ADDR must be a host:port address: %w", err)
	}
	return cfg, nil
}

func valueOrDefault(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && strings.TrimSpace(value) != "" {
		return strings.TrimSpace(value)
	}
	return fallback
}

func required(key string) string {
	return strings.TrimSpace(os.Getenv(key))
}
