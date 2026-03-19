package config

import (
	"strings"
	"testing"
)

func TestLoadUsesDefaults(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://example")
	t.Setenv("JWT_SECRET", "super-secret")
	t.Setenv("HTTP_ADDR", "")
	t.Setenv("JWT_TTL_MIN", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if cfg.HTTPAddr != ":8080" {
		t.Fatalf("expected default HTTP addr :8080, got %q", cfg.HTTPAddr)
	}
	if cfg.JWTTTLMin != 120 {
		t.Fatalf("expected default ttl 120, got %d", cfg.JWTTTLMin)
	}
}

func TestLoadUsesCustomHTTPAddr(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://example")
	t.Setenv("JWT_SECRET", "super-secret")
	t.Setenv("HTTP_ADDR", ":9090")
	t.Setenv("JWT_TTL_MIN", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if cfg.HTTPAddr != ":9090" {
		t.Fatalf("expected custom HTTP addr :9090, got %q", cfg.HTTPAddr)
	}
}

func TestLoadRejectsNonPositiveTTL(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://example")
	t.Setenv("JWT_SECRET", "super-secret")
	t.Setenv("JWT_TTL_MIN", "0")

	_, err := Load()
	if err == nil {
		t.Fatalf("expected Load to fail for non-positive TTL")
	}
	if !strings.Contains(err.Error(), "JWT_TTL_MIN") {
		t.Fatalf("expected JWT_TTL_MIN error, got: %v", err)
	}
}

func TestLoadRejectsInvalidTTL(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://example")
	t.Setenv("JWT_SECRET", "super-secret")
	t.Setenv("JWT_TTL_MIN", "abc")

	_, err := Load()
	if err == nil {
		t.Fatalf("expected Load to fail for invalid TTL")
	}
	if !strings.Contains(err.Error(), "invalid JWT_TTL_MIN") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoadRequiresDatabaseURL(t *testing.T) {
	t.Setenv("DATABASE_URL", "")
	t.Setenv("JWT_SECRET", "super-secret")
	t.Setenv("JWT_TTL_MIN", "")

	_, err := Load()
	if err == nil {
		t.Fatalf("expected Load to fail when DATABASE_URL is missing")
	}
	if !strings.Contains(err.Error(), "DATABASE_URL") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoadRequiresJWTSecret(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://example")
	t.Setenv("JWT_SECRET", "")
	t.Setenv("JWT_TTL_MIN", "")

	_, err := Load()
	if err == nil {
		t.Fatalf("expected Load to fail when JWT_SECRET is missing")
	}
	if !strings.Contains(err.Error(), "JWT_SECRET") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoadPreservesProvidedDatabaseURL(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://user:pass@db:5432/chess?sslmode=disable")
	t.Setenv("JWT_SECRET", "super-secret")
	t.Setenv("JWT_TTL_MIN", "60")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if !strings.Contains(cfg.DatabaseURL, "db:5432/chess") {
		t.Fatalf("expected database URL to be preserved, got %q", cfg.DatabaseURL)
	}
	if cfg.JWTTTLMin != 60 {
		t.Fatalf("expected ttl 60, got %d", cfg.JWTTTLMin)
	}
}
