package domain

import (
	"testing"
)

func TestNewMeetup(t *testing.T) {

	tests := []struct {
		name        string
		title       string
		capacity    int
		expectError bool
	}{
		{"Valid meetup", "Go Architecture", 100, false},
		{"Empty title", "", 100, true},
		{"Zero capacity", "Go Architecture", 0, true},
		{"Negative capacity", "Go Architecture", -50, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			m, err := NewMeetup(tt.title, tt.capacity)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected an error for invalid input, but got none")
				}
				if m != nil {
					t.Errorf("expected meetup to be nil on error, but got object")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if m == nil {
					t.Errorf("expected meetup to be created, got nil")
				}
				if m != nil && m.Status != StatusDraft {
					t.Errorf("expected initial status to be %q, got %q", StatusDraft, m.Status)
				}
			}
		})
	}
}

func TestMeetup_ChangeStatus(t *testing.T) {
	tests := []struct {
		name        string
		startStatus MeetupStatus
		newStatus   MeetupStatus
		expectError bool
	}{
		{"Draft -> Published (Valid)", StatusDraft, StatusPublished, false},
		{"Draft -> Canceled (Valid)", StatusDraft, StatusCanceled, false},
		{"Published -> Archived (Valid)", StatusPublished, StatusArchived, false},
		{"Published -> Canceled (Valid)", StatusPublished, StatusCanceled, false},
		{"Draft -> Archived (Invalid)", StatusDraft, StatusArchived, true},
		{"Published -> Draft (Invalid)", StatusPublished, StatusDraft, true},
		{"Archived -> Published (Invalid)", StatusArchived, StatusPublished, true},
		{"Archived -> Draft (Invalid)", StatusArchived, StatusDraft, true},
		{"Canceled -> Draft (Invalid)", StatusCanceled, StatusDraft, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			m := &Meetup{
				Status: tt.startStatus,
			}

			err := m.ChangeStatus(tt.newStatus)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error for transition %s -> %s, got none", tt.startStatus, tt.newStatus)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for transition %s -> %s: %v", tt.startStatus, tt.newStatus, err)
				}
				if m.Status != tt.newStatus {
					t.Errorf("expected status to change to %s, got %s", tt.newStatus, m.Status)
				}
			}
		})
	}
}
