package handlers

import (
	"net/http"
	"testing"

	"chess-training/internal/service"
)

func TestRegisterErrorResponse(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantCode   string
	}{
		{
			name:       "username taken",
			err:        service.ErrUsernameTaken,
			wantStatus: http.StatusConflict,
			wantCode:   "USERNAME_TAKEN",
		},
		{
			name:       "email taken",
			err:        service.ErrEmailTaken,
			wantStatus: http.StatusConflict,
			wantCode:   "EMAIL_TAKEN",
		},
		{
			name:       "invalid username",
			err:        service.ErrInvalidUsername,
			wantStatus: http.StatusBadRequest,
			wantCode:   "INVALID_USERNAME",
		},
		{
			name:       "weak password",
			err:        service.ErrWeakPassword,
			wantStatus: http.StatusBadRequest,
			wantCode:   "WEAK_PASSWORD",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStatus, gotCode, _ := registerErrorResponse(tt.err)
			if gotStatus != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, gotStatus)
			}
			if gotCode != tt.wantCode {
				t.Fatalf("expected code %q, got %q", tt.wantCode, gotCode)
			}
		})
	}
}
