package create_user

import (
	"context"
	"errors"
	"testing"

	"github.com/anechytailenko/Microservices_practice04/internal/shared"
	"github.com/anechytailenko/Microservices_practice04/internal/users/domain"
)

type fakeRepo struct {
	savedUser *domain.User
	saveErr   error
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{}
}

func (f *fakeRepo) Save(ctx context.Context, user *domain.User) error {
	if f.saveErr != nil {
		return f.saveErr
	}
	f.savedUser = user
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
			name: "Success: User Created and Saved",
			command: Command{
				FirstName: "Alex",
				LastName:  "Johnson",
				Email:     "alex.j@example.com",
			},
			setupMock:     func(repo *fakeRepo) {},
			expectedError: nil,
		},
		{
			name: "Error: Domain Rule Violation (Empty First Name)",
			command: Command{
				FirstName: "",
				LastName:  "Johnson",
				Email:     "alex.j@example.com",
			},
			setupMock:     func(repo *fakeRepo) {},
			expectedError: shared.Error{Type: shared.ErrorTypeValidation},
		},
		{
			name: "Error: Domain Rule Violation (Invalid Email)",
			command: Command{
				FirstName: "Alex",
				LastName:  "Johnson",
				Email:     "not-an-email",
			},
			setupMock:     func(repo *fakeRepo) {},
			expectedError: shared.Error{Type: shared.ErrorTypeValidation},
		},
		{
			name: "Error: Database Save Fails",
			command: Command{
				FirstName: "Alex",
				LastName:  "Johnson",
				Email:     "alex.j@example.com",
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
					t.Errorf("expected valid user ID, got empty string")
				}

				if repo.savedUser == nil {
					t.Fatalf("expected user to be passed to repo.Save, but it was nil")
				}

				if repo.savedUser.FirstName != tt.command.FirstName {
					t.Errorf("expected saved first name %q, got %q", tt.command.FirstName, repo.savedUser.FirstName)
				}
				if repo.savedUser.LastName != tt.command.LastName {
					t.Errorf("expected saved last name %q, got %q", tt.command.LastName, repo.savedUser.LastName)
				}
				if repo.savedUser.Email != tt.command.Email {
					t.Errorf("expected saved email %q, got %q", tt.command.Email, repo.savedUser.Email)
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
