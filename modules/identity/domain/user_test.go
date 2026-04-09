package domain

import (
	"testing"
	"time"
)

func TestNewUser(t *testing.T) {
	now := time.Date(2026, 4, 9, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		id      string
		user    string
		email   string
		wantErr error
	}{
		{
			name:  "creates valid user",
			id:    "user-1",
			user:  "Ozgur",
			email: "ozgur@example.com",
		},
		{
			name:    "returns error when name is empty",
			id:      "user-1",
			user:    "",
			email:   "ozgur@example.com",
			wantErr: ErrInvalidUserName,
		},
		{
			name:    "returns error when email is empty",
			id:      "user-1",
			user:    "Ozgur",
			email:   "",
			wantErr: ErrInvalidUserEmail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewUser(tt.id, tt.user, tt.email, now)
			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if got.ID != tt.id {
				t.Fatalf("expected id %s, got %s", tt.id, got.ID)
			}

			if got.Name != tt.user {
				t.Fatalf("expected name %s, got %s", tt.user, got.Name)
			}

			if got.Email != tt.email {
				t.Fatalf("expected email %s, got %s", tt.email, got.Email)
			}

			if got.CreatedAt != now {
				t.Fatalf("expected created at %v, got %v", now, got.CreatedAt)
			}
		})
	}
}
