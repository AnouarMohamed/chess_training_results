package middleware

import "testing"

func TestBearerTokenFromAuthorizationHeader(t *testing.T) {
	tests := []struct {
		name    string
		header  string
		want    string
		wantErr bool
	}{
		{name: "valid bearer", header: "Bearer abc.def.ghi", want: "abc.def.ghi"},
		{name: "case insensitive scheme", header: "bearer token123", want: "token123"},
		{name: "extra spaces", header: "Bearer   token-xyz", want: "token-xyz"},
		{name: "missing header", header: "", wantErr: true},
		{name: "missing token", header: "Bearer", wantErr: true},
		{name: "wrong scheme", header: "Basic 123", wantErr: true},
		{name: "too many parts", header: "Bearer first second", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := bearerTokenFromAuthorizationHeader(tt.header)
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
				t.Fatalf("expected token %q, got %q", tt.want, got)
			}
		})
	}
}
