package handlers

import (
	db "chess-training/internal/db/sqlc"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestNewPlayerResponseMapsSafeFields(t *testing.T) {
	id := uuid.MustParse("019f2cd9-6846-76f0-a7fe-8f0c8ce939e2")
	player := db.Player{
		ID:           pgtype.UUID{Bytes: [16]byte(id), Valid: true},
		Username:     "anouar",
		Email:        pgtype.Text{String: "anouar@example.com", Valid: true},
		PasswordHash: "should-not-leak",
	}

	response := newPlayerResponse(player)
	body, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("failed to marshal response: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("failed to unmarshal response JSON: %v", err)
	}

	if payload["id"] != id.String() {
		t.Fatalf("expected id %q, got %v", id.String(), payload["id"])
	}
	if payload["username"] != "anouar" {
		t.Fatalf("expected username anouar, got %v", payload["username"])
	}
	if payload["email"] != "anouar@example.com" {
		t.Fatalf("expected email anouar@example.com, got %v", payload["email"])
	}
	if _, ok := payload["passwordHash"]; ok {
		t.Fatalf("response unexpectedly exposes passwordHash")
	}
}

func TestNewPlayerResponseOmitsInvalidEmail(t *testing.T) {
	player := db.Player{
		ID:       pgtype.UUID{Valid: false},
		Username: "reader",
		Email:    pgtype.Text{Valid: false},
	}

	response := newPlayerResponse(player)
	body, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("failed to marshal response: %v", err)
	}

	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("failed to unmarshal response JSON: %v", err)
	}

	if payload["id"] != "" {
		t.Fatalf("expected empty id for invalid UUID, got %v", payload["id"])
	}
	if _, ok := payload["email"]; ok {
		t.Fatalf("expected email field to be omitted when invalid")
	}
}
