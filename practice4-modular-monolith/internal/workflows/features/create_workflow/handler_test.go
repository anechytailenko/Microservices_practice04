package create_workflow

import (
	"context"
	"errors"
	"testing"

	"github.com/anechytailenko/Microservices_practice04/internal/shared"
	"github.com/anechytailenko/Microservices_practice04/internal/workflows/domain"
)

type fakeRepo struct {
	saveErr error
	savedWf *domain.Workflow
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{}
}

func (f *fakeRepo) Save(ctx context.Context, workflow *domain.Workflow, eventID string, eventType string, eventPayload []byte) error {
	if f.saveErr != nil {
		return f.saveErr
	}
	f.savedWf = workflow
	return nil
}

func TestHandler_Handle(t *testing.T) {
	tests := []struct {
		name          string
		cmd           Command
		setupMock     func(repo *fakeRepo)
		expectedError error
	}{
		{
			name: "Success: Workflow Created",
			cmd: Command{
				UserID:   "user-1",
				MeetupID: "meetup-1",
			},
			setupMock:     func(repo *fakeRepo) {},
			expectedError: nil,
		},
		{
			name: "Error: Missing IDs",
			cmd: Command{
				UserID:   "",
				MeetupID: "",
			},
			setupMock:     func(repo *fakeRepo) {},
			expectedError: shared.Error{Type: shared.ErrorTypeValidation},
		},
		{
			name: "Error: DB Save Fails",
			cmd: Command{
				UserID:   "user-1",
				MeetupID: "meetup-1",
			},
			setupMock: func(repo *fakeRepo) {
				repo.saveErr = errors.New("db error")
			},
			expectedError: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newFakeRepo()
			tt.setupMock(repo)
			handler := NewHandler(repo)

			dto, err := handler.Handle(context.Background(), tt.cmd)

			if tt.expectedError == nil {
				if err != nil {
					t.Fatalf("expected no error, got: %v", err)
				}
				if dto == nil {
					t.Fatalf("expected valid DTO, got nil")
				}
				if dto.ID == "" {
					t.Errorf("expected generated ID")
				}
				if dto.Type != "join-meetup" {
					t.Errorf("expected type 'join-meetup', got %q", dto.Type)
				}
				if dto.State != string(domain.StateInitializing) {
					t.Errorf("expected state %q, got %q", domain.StateInitializing, dto.State)
				}

				if repo.savedWf == nil {
					t.Errorf("expected workflow to be passed to repo.Save")
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
