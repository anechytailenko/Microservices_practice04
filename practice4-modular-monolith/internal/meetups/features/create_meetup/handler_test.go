package create_meetup

import (
	"context"
	"errors"
	"testing"

	"github.com/anechytailenko/Microservices_practice04/internal/meetups/domain"
	"github.com/anechytailenko/Microservices_practice04/internal/shared"
)

type fakeRepo struct {
	savedMeetup *domain.Meetup
	saveErr     error
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{}
}

func (f *fakeRepo) Save(ctx context.Context, meetup *domain.Meetup) error {
	if f.saveErr != nil {
		return f.saveErr
	}
	f.savedMeetup = meetup
	return nil
}

func TestHandler_Handle(t *testing.T) {
	tests := []struct {
		name          string
		command       Command
		setupMock     func(repo *fakeRepo)
		expectedError error
	}{
		{
			name: "Success: Meetup Created and Saved",
			command: Command{
				Title:    "Go Architecture Meetup",
				Capacity: 100,
			},
			setupMock:     func(repo *fakeRepo) {},
			expectedError: nil,
		},
		{
			name: "Error: Domain Rule Violation (Empty Title)",
			command: Command{
				Title:    "",
				Capacity: 100,
			},
			setupMock:     func(repo *fakeRepo) {},
			expectedError: shared.Error{Type: shared.ErrorTypeValidation},
		},
		{
			name: "Error: Domain Rule Violation (Negative Capacity)",
			command: Command{
				Title:    "Valid Title",
				Capacity: -5,
			},
			setupMock:     func(repo *fakeRepo) {},
			expectedError: shared.Error{Type: shared.ErrorTypeValidation},
		},
		{
			name: "Error: Database Save Fails",
			command: Command{
				Title:    "Valid Title",
				Capacity: 100,
			},
			setupMock: func(repo *fakeRepo) {
				repo.saveErr = errors.New("database timeout")
			},
			expectedError: errors.New("database timeout"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newFakeRepo()
			tt.setupMock(repo)
			handler := NewHandler(repo)

			id, err := handler.Handle(context.Background(), tt.command)

			if tt.expectedError == nil {
				if err != nil {
					t.Errorf("expected no error, got: %v", err)
				}

				if id == "" {
					t.Errorf("expected valid meetup ID, got empty string")
				}

				if repo.savedMeetup == nil {
					t.Fatalf("expected meetup to be passed to repo.Save, but it was nil")
				}

				if repo.savedMeetup.Title != tt.command.Title {
					t.Errorf("expected saved title %q, got %q", tt.command.Title, repo.savedMeetup.Title)
				}
				if repo.savedMeetup.Capacity != tt.command.Capacity {
					t.Errorf("expected saved capacity %d, got %d", tt.command.Capacity, repo.savedMeetup.Capacity)
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
