package util

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
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

func TestParseJWTRejectsUnexpectedSigningMethod(t *testing.T) {
	secret := "test-secret"
	claims := Claims{
		Sub:      "user-123",
		Username: "anouar",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("failed to sign test token: %v", err)
	}

	if _, err := ParseJWT(secret, tokenStr); err == nil {
		t.Fatalf("expected ParseJWT to reject non-HS256 token")
	}
}
