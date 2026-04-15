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
				if u.Meetups == nil {
					t.Errorf("expected meetups slice to be initialized, got nil")
				}
				if len(u.Meetups) != 0 {
					t.Errorf("expected meetups slice to be empty, got length %d", len(u.Meetups))
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

func TestUser_AddMeetup(t *testing.T) {
	u := &User{
		Meetups: make([]string, 0),
	}

	meetupID := "meetup-123"

	u.AddMeetup(meetupID)
	if len(u.Meetups) != 1 || u.Meetups[0] != meetupID {
		t.Errorf("expected meetup %q to be added, got %v", meetupID, u.Meetups)
	}

	u.AddMeetup(meetupID)
	if len(u.Meetups) != 1 {
		t.Errorf("expected idempotency to prevent duplicates, expected length 1, got %d", len(u.Meetups))
	}

	u.AddMeetup("meetup-456")
	if len(u.Meetups) != 2 {
		t.Errorf("expected length 2 after adding a different meetup, got %d", len(u.Meetups))
	}
}

func TestUser_RemoveMeetup(t *testing.T) {
	u := &User{
		Meetups: []string{"meetup-123", "meetup-456"},
	}

	u.RemoveMeetup("meetup-123")
	if len(u.Meetups) != 1 || u.Meetups[0] != "meetup-456" {
		t.Errorf("expected 'meetup-456' to remain, got %v", u.Meetups)
	}

	u.RemoveMeetup("meetup-999")
	if len(u.Meetups) != 1 {
		t.Errorf("expected idempotency to handle non-existent meetups, expected length 1, got %d", len(u.Meetups))
	}

	u.RemoveMeetup("meetup-456")
	if len(u.Meetups) != 0 {
		t.Errorf("expected meetups to be empty, got length %d", len(u.Meetups))
	}
}
