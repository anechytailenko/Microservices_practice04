package change_status

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/anechytailenko/Microservices_practice04/internal/meetups/domain"
	"github.com/anechytailenko/Microservices_practice04/internal/shared"
)

type fakeRepo struct {
	meetups   map[domain.MeetupID]*domain.Meetup
	getErr    error
	updateErr error
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{
		meetups: make(map[domain.MeetupID]*domain.Meetup),
	}
}

func (f *fakeRepo) GetByID(ctx context.Context, id domain.MeetupID) (*domain.Meetup, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	return f.meetups[id], nil
}

func (f *fakeRepo) Update(ctx context.Context, meetup *domain.Meetup) error {
	if f.updateErr != nil {
		return f.updateErr
	}
	f.meetups[meetup.ID] = meetup
	return nil
}

func TestHandler_Handle(t *testing.T) {
	meetupID := domain.MeetupID("12345")

	initialMeetup := domain.Meetup{
		ID:          meetupID,
		Title:       "Test Meetup",
		Capacity:    100,
		OwnerUserID: "user-123",
		Status:      domain.StatusDraft,
		CreatedAt:   time.Now(),
	}

	tests := []struct {
		name          string
		command       Command
		setupMock     func(repo *fakeRepo)
		expectedError error
	}{
		{
			name:    "Success: Draft -> Published",
			command: Command{MeetupID: string(meetupID), Status: string(domain.StatusPublished)},
			setupMock: func(repo *fakeRepo) {
				m := initialMeetup
				repo.meetups[meetupID] = &m
			},
			expectedError: nil,
		},
		{
			name:          "Error: Meetup Not Found",
			command:       Command{MeetupID: "99999", Status: string(domain.StatusPublished)},
			setupMock:     func(repo *fakeRepo) {},
			expectedError: shared.Error{Type: shared.ErrorTypeNotFound},
		},
		{
			name:    "Error: Domain Rule Violation (Draft -> Archived)",
			command: Command{MeetupID: string(meetupID), Status: string(domain.StatusArchived)},
			setupMock: func(repo *fakeRepo) {
				m := initialMeetup
				repo.meetups[meetupID] = &m
			},
			expectedError: shared.Error{Type: shared.ErrorTypeConflict},
		},
		{
			name:    "Error: Database GetByID fails",
			command: Command{MeetupID: string(meetupID), Status: string(domain.StatusPublished)},
			setupMock: func(repo *fakeRepo) {
				repo.getErr = errors.New("database connection lost")
			},
			expectedError: errors.New("database connection lost"),
		},
		{
			name:    "Error: Database Update fails",
			command: Command{MeetupID: string(meetupID), Status: string(domain.StatusPublished)},
			setupMock: func(repo *fakeRepo) {
				m := initialMeetup
				repo.meetups[meetupID] = &m
				repo.updateErr = errors.New("failed to write to database")
			},
			expectedError: errors.New("failed to write to database"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newFakeRepo()
			tt.setupMock(repo)
			handler := NewHandler(repo)

			err := handler.Handle(context.Background(), tt.command)

			if tt.expectedError == nil {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if updatedMeetup := repo.meetups[meetupID]; updatedMeetup.Status != domain.MeetupStatus(tt.command.Status) {
					t.Errorf("expected status %s, got %s", tt.command.Status, updatedMeetup.Status)
				}
			} else {
				if err == nil {
					t.Errorf("expected error, got nil")
					return
				}

				if expectedSharedErr, ok := tt.expectedError.(shared.Error); ok {
					var actualSharedErr shared.Error

					if errors.As(err, &actualSharedErr) {
						if actualSharedErr.Type != expectedSharedErr.Type {
							t.Errorf("expected error type %s, got %s", expectedSharedErr.Type, actualSharedErr.Type)
						}
					} else {
						t.Errorf("expected shared.Error, got generic error: %v", err)
					}
				}
			}
		})
	}
}
