package create_meetup

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/anechytailenko/Microservices_practice04/internal/meetups/domain"
	"github.com/anechytailenko/Microservices_practice04/internal/shared"
	"github.com/anechytailenko/Microservices_practice04/internal/shared/ctxutil"
	"github.com/anechytailenko/Microservices_practice04/internal/shared/events"
)

type fakeRepo struct {
	savedMeetup    *domain.Meetup
	savedEventID   string
	savedEventType string
	savedEvent     events.MeetupCreatedEvent
	savedPayload   []byte
	saveErr        error
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{}
}

func (f *fakeRepo) Save(ctx context.Context, meetup *domain.Meetup, eventID string, eventType string, payload []byte) error {
	if f.saveErr != nil {
		return f.saveErr
	}
	f.savedMeetup = meetup
	f.savedEventID = eventID
	f.savedEventType = eventType
	f.savedPayload = payload

	var evt events.MeetupCreatedEvent
	if err := json.Unmarshal(payload, &evt); err == nil {
		f.savedEvent = evt
	}

	return nil
}

type fakeUserValidator struct {
	validateErr error
}

func newFakeUserValidator() *fakeUserValidator {
	return &fakeUserValidator{}
}

func (f *fakeUserValidator) ValidateUserExists(ctx context.Context, userID string) error {
	return f.validateErr
}

func TestHandler_Handle(t *testing.T) {
	tests := []struct {
		name          string
		command       Command
		setupMock     func(repo *fakeRepo, userValidator *fakeUserValidator)
		expectedError error
	}{
		{
			name: "Success: Meetup Created and Saved",
			command: Command{
				Title:       "Go Architecture Meetup",
				Capacity:    100,
				OwnerUserID: "user-123",
			},
			setupMock:     func(repo *fakeRepo, userValidator *fakeUserValidator) {},
			expectedError: nil,
		},
		{
			name: "Error: User Validation Fails (Not Found)",
			command: Command{
				Title:       "Go Architecture Meetup",
				Capacity:    100,
				OwnerUserID: "non-existent-user",
			},
			setupMock: func(repo *fakeRepo, userValidator *fakeUserValidator) {
				userValidator.validateErr = shared.NewValidationError("owner_user_id 'non-existent-user' does not exist")
			},
			expectedError: shared.Error{Type: shared.ErrorTypeValidation},
		},
		{
			name: "Error: Users Service Unreachable",
			command: Command{
				Title:       "Go Architecture Meetup",
				Capacity:    100,
				OwnerUserID: "user-123",
			},
			setupMock: func(repo *fakeRepo, userValidator *fakeUserValidator) {
				userValidator.validateErr = shared.NewServiceUnavailableError("users service is unreachable")
			},
			expectedError: shared.Error{Type: shared.ErrorTypeServiceUnavailable},
		},
		{
			name: "Error: Domain Rule Violation (Empty Title)",
			command: Command{
				Title:       "",
				Capacity:    100,
				OwnerUserID: "user-123",
			},
			setupMock:     func(repo *fakeRepo, userValidator *fakeUserValidator) {},
			expectedError: shared.Error{Type: shared.ErrorTypeValidation},
		},
		{
			name: "Error: Domain Rule Violation (Negative Capacity)",
			command: Command{
				Title:       "Valid Title",
				Capacity:    -5,
				OwnerUserID: "user-123",
			},
			setupMock:     func(repo *fakeRepo, userValidator *fakeUserValidator) {},
			expectedError: shared.Error{Type: shared.ErrorTypeValidation},
		},
		{
			name: "Error: Database Save Fails",
			command: Command{
				Title:       "Valid Title",
				Capacity:    100,
				OwnerUserID: "user-123",
			},
			setupMock: func(repo *fakeRepo, userValidator *fakeUserValidator) {
				repo.saveErr = errors.New("database timeout")
			},
			expectedError: errors.New("database timeout"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newFakeRepo()
			userValidator := newFakeUserValidator()
			tt.setupMock(repo, userValidator)

			handler := NewHandler(repo, userValidator)

			expectedCorrID := "test-correlation-id-999"
			ctx := ctxutil.WithCorrelationID(context.Background(), expectedCorrID)

			id, err := handler.Handle(ctx, tt.command)

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
				if repo.savedMeetup.OwnerUserID != tt.command.OwnerUserID {
					t.Errorf("expected saved owner user ID %q, got %q", tt.command.OwnerUserID, repo.savedMeetup.OwnerUserID)
				}
				if repo.savedEventType != "meetup.created" {
					t.Errorf("expected event type 'meetup.created', got %q", repo.savedEventType)
				}
				if repo.savedEvent.CorrelationID != expectedCorrID {
					t.Errorf("expected event CorrelationID %q, got %q", expectedCorrID, repo.savedEvent.CorrelationID)
				}
				if len(repo.savedPayload) == 0 {
					t.Errorf("expected non-empty JSON payload, got empty")
				}
				if repo.savedEventID == "" {
					t.Errorf("expected non-empty EventID to be passed to repo")
				}
				if repo.savedEventID != repo.savedEvent.EventID {
					t.Errorf("expected EventID passed to DB (%q) to match EventID inside JSON payload (%q)", repo.savedEventID, repo.savedEvent.EventID)
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
