package domain

import (
	"strings"
	"testing"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name        string
		firstName   string
		lastName    string
		email       string
		expectError bool
	}{
		{"Valid User", "Alex", "Johnson", "alex.j@example.com", false},
		{"Valid User (requires trimming)", "  Alex  ", "  Johnson  ", "  alex.j@example.com  ", false},
		{"Empty First Name", "", "Johnson", "alex.j@example.com", true},
		{"Empty Last Name", "Alex", "", "alex.j@example.com", true},
		{"Empty Email", "Alex", "Johnson", "", true},
		{"Invalid Email Format", "Alex", "Johnson", "invalid-email-format", true},
		{"Spaces only in First Name", "   ", "Johnson", "alex.j@example.com", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := NewUser(tt.firstName, tt.lastName, tt.email)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected validation error for input, but got none")
				}
				if u != nil {
					t.Errorf("expected user to be nil on error, got: %v", u)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if u == nil {
					t.Fatalf("expected user to be created, got nil")
				}

				expectedFirst := strings.TrimSpace(tt.firstName)
				expectedLast := strings.TrimSpace(tt.lastName)
				expectedEmail := strings.TrimSpace(tt.email)

				if u.FirstName != expectedFirst {
					t.Errorf("expected first name %q, got %q", expectedFirst, u.FirstName)
				}
				if u.LastName != expectedLast {
					t.Errorf("expected last name %q, got %q", expectedLast, u.LastName)
				}
				if u.Email != expectedEmail {
					t.Errorf("expected email %q, got %q", expectedEmail, u.Email)
				}

				if u.ID == "" {
					t.Errorf("expected user ID to be generated, got empty string")
				}
			}
		})
	}
}

func TestUser_DisplayName(t *testing.T) {
	u := &User{
		FirstName: "Alex",
		LastName:  "Johnson",
	}

	expected := "Alex Johnson"
	actual := u.DisplayName()

	if actual != expected {
		t.Errorf("expected DisplayName %q, got %q", expected, actual)
	}
}
