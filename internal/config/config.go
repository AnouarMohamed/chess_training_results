package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	HTTPAddr    string
	DatabaseURL string
	JWTSecret   string
	JWTTTLMin   int
	AppEnv      string
}

func Load() (Config, error) {
	ttl := 120
	if v := os.Getenv("JWT_TTL_MIN"); v != "" {
		i, err := strconv.Atoi(v)
		if err != nil {
			return Config{}, fmt.Errorf("invalid JWT_TTL_MIN: %w", err)
		}
		if i <= 0 {
			return Config{}, fmt.Errorf("invalid JWT_TTL_MIN: must be > 0")
		}
		ttl = i
	}
	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	appEnv, err := resolveAppEnv(os.Getenv("APP_ENV"))
	if err != nil {
		return Config{}, err
	}

	databaseURL, err := requiredEnv("DATABASE_URL")
	if err != nil {
		return Config{}, err
	}
	jwtSecret, err := requiredEnv("JWT_SECRET")
	if err != nil {
		return Config{}, err
	}

	return Config{
		HTTPAddr:    addr,
		DatabaseURL: databaseURL,
		JWTSecret:   jwtSecret,
		JWTTTLMin:   ttl,
		AppEnv:      appEnv,
	}, nil
}

func requiredEnv(k string) (string, error) {
	v := os.Getenv(k)
	if v == "" {
		return "", fmt.Errorf("missing env var: %s", k)
	}
	return v, nil
}

func resolveAppEnv(raw string) (string, error) {
	if raw == "" {
		return "dev", nil
	}
	switch raw {
	case "dev", "test", "staging", "prod":
		return raw, nil
	default:
		return "", fmt.Errorf("invalid APP_ENV: %s", raw)
	}
}
