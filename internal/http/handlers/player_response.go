package handlers

import (
	db "chess-training/internal/db/sqlc"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type playerResponse struct {
	ID       string  `json:"id"`
	Username string  `json:"username"`
	Email    *string `json:"email,omitempty"`
}

func newPlayerResponse(player db.Player) playerResponse {
	return playerResponse{
		ID:       uuidToString(player.ID),
		Username: player.Username,
		Email:    nullableText(player.Email),
	}
}

func uuidToString(value pgtype.UUID) string {
	if !value.Valid {
		return ""
	}
	parsed, err := uuid.FromBytes(value.Bytes[:])
	if err != nil {
		return ""
	}
	return parsed.String()
}

func nullableText(value pgtype.Text) *string {
	if !value.Valid || value.String == "" {
		return nil
	}
	v := value.String
	return &v
}
