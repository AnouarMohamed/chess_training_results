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
	t.Setenv("APP_ENV", "")

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
	if cfg.AppEnv != "dev" {
		t.Fatalf("expected default app env dev, got %q", cfg.AppEnv)
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

func TestLoadRejectsInvalidAppEnv(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://example")
	t.Setenv("JWT_SECRET", "super-secret")
	t.Setenv("JWT_TTL_MIN", "")
	t.Setenv("APP_ENV", "production")

	_, err := Load()
	if err == nil {
		t.Fatalf("expected Load to fail when APP_ENV is invalid")
	}
	if !strings.Contains(err.Error(), "APP_ENV") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestLoadAcceptsSupportedAppEnv(t *testing.T) {
	supported := []string{"dev", "test", "staging", "prod"}
	for _, env := range supported {
		t.Run(env, func(t *testing.T) {
			t.Setenv("DATABASE_URL", "postgres://example")
			t.Setenv("JWT_SECRET", "super-secret")
			t.Setenv("JWT_TTL_MIN", "")
			t.Setenv("APP_ENV", env)

			cfg, err := Load()
			if err != nil {
				t.Fatalf("Load returned error: %v", err)
			}
			if cfg.AppEnv != env {
				t.Fatalf("expected app env %q, got %q", env, cfg.AppEnv)
			}
		})
	}
}
