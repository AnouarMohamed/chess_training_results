package service

import "testing"

func TestNormalizeUsername(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "valid with underscore", input: "anouar_dev", want: "anouar_dev"},
		{name: "valid with spaces around", input: "  chessUser123  ", want: "chessUser123"},
		{name: "too short", input: "ab", wantErr: true},
		{name: "too long", input: "abcdefghijklmnopqrstuvwxyz1234567", wantErr: true},
		{name: "invalid char dash", input: "anouar-dev", wantErr: true},
		{name: "invalid char space inside", input: "anouar dev", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalizeUsername(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{name: "valid", input: "superStrong1"},
		{name: "trims spaces still valid", input: "  strongpass  "},
		{name: "too short", input: "abc123", wantErr: true},
		{name: "blank", input: "   ", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePassword(tt.input)
			if tt.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
