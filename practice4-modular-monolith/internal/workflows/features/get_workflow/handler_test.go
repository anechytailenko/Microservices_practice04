package get_workflow

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/anechytailenko/Microservices_practice04/internal/shared"
	"github.com/anechytailenko/Microservices_practice04/internal/workflows/domain"
)

type fakeRepo struct {
	workflows map[domain.WorkflowID]*domain.Workflow
	getErr    error
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{
		workflows: make(map[domain.WorkflowID]*domain.Workflow),
	}
}

func (f *fakeRepo) GetByID(ctx context.Context, id domain.WorkflowID) (*domain.Workflow, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	return f.workflows[id], nil
}

func TestHandler_Handle(t *testing.T) {
	wfID := domain.WorkflowID("wf-123")
	now := time.Now().UTC()
	errReason := "something went wrong"

	initialWorkflow := domain.Workflow{
		ID:        wfID,
		Type:      "join-meetup",
		State:     domain.StateFailed,
		LastError: &errReason,
		CreatedAt: now,
		UpdatedAt: now,
	}

	tests := []struct {
		name          string
		query         Query
		setupMock     func(repo *fakeRepo)
		expectedError error
	}{
		{
			name:  "Success: Workflow Found",
			query: Query{WorkflowID: string(wfID)},
			setupMock: func(repo *fakeRepo) {
				w := initialWorkflow
				repo.workflows[wfID] = &w
			},
			expectedError: nil,
		},
		{
			name:          "Error: Empty ID",
			query:         Query{WorkflowID: ""},
			setupMock:     func(repo *fakeRepo) {},
			expectedError: shared.Error{Type: shared.ErrorTypeValidation},
		},
		{
			name:          "Error: Workflow Not Found",
			query:         Query{WorkflowID: "non-existent"},
			setupMock:     func(repo *fakeRepo) {},
			expectedError: shared.Error{Type: shared.ErrorTypeNotFound},
		},
		{
			name:  "Error: DB Error",
			query: Query{WorkflowID: string(wfID)},
			setupMock: func(repo *fakeRepo) {
				repo.getErr = errors.New("db connection lost")
			},
			expectedError: errors.New("db connection lost"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newFakeRepo()
			tt.setupMock(repo)
			handler := NewHandler(repo)

			dto, err := handler.Handle(context.Background(), tt.query)

			if tt.expectedError == nil {
				if err != nil {
					t.Fatalf("expected no error, got: %v", err)
				}
				if dto == nil {
					t.Fatalf("expected valid DTO, got nil")
				}
				if dto.ID != string(initialWorkflow.ID) {
					t.Errorf("expected ID %q, got %q", initialWorkflow.ID, dto.ID)
				}
				if dto.State != string(initialWorkflow.State) {
					t.Errorf("expected State %q, got %q", initialWorkflow.State, dto.State)
				}
				if dto.LastError == nil || *dto.LastError != errReason {
					t.Errorf("expected LastError %q", errReason)
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
