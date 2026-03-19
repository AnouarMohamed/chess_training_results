package util

import (
	"testing"
)

func TestSignAndParseJWT(t *testing.T) {
	secret := "test-secret"
	userID := "3f9d4efd-dbb6-46f8-a13f-b46b2d6b4d2f"
	username := "anouar"

	token, err := SignJWT(secret, userID, username, 60)
	if err != nil {
		t.Fatalf("SignJWT returned error: %v", err)
	}

	claims, err := ParseJWT(secret, token)
	if err != nil {
		t.Fatalf("ParseJWT returned error: %v", err)
	}

	if claims.Sub != userID {
		t.Fatalf("expected subject %q, got %q", userID, claims.Sub)
	}
	if claims.Username != username {
		t.Fatalf("expected username %q, got %q", username, claims.Username)
	}
	if claims.ExpiresAt == nil || claims.IssuedAt == nil {
		t.Fatalf("expected issued and expiry claims to be set")
	}
}

func TestParseJWTRejectsWrongSecret(t *testing.T) {
	token, err := SignJWT("secret-a", "uid", "name", 60)
	if err != nil {
		t.Fatalf("SignJWT returned error: %v", err)
	}

	if _, err := ParseJWT("secret-b", token); err == nil {
		t.Fatalf("expected ParseJWT to fail with wrong secret")
	}
}
