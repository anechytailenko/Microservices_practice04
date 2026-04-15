package get_user

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/anechytailenko/Microservices_practice04/internal/shared"
	"github.com/anechytailenko/Microservices_practice04/internal/users/domain"
)

type fakeRepo struct {
	users  map[domain.UserID]*domain.User
	getErr error
}

func newFakeRepo() *fakeRepo {
	return &fakeRepo{
		users: make(map[domain.UserID]*domain.User),
	}
}

func (f *fakeRepo) GetByID(ctx context.Context, id domain.UserID) (*domain.User, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}

	return f.users[id], nil
}

func TestHandler_Handle(t *testing.T) {
	userID := domain.UserID("user-123")

	initialUser := domain.User{
		ID:        userID,
		FirstName: "Alex",
		LastName:  "Johnson",
		Email:     "alex.j@example.com",
		Meetups:   []string{"meetup-1", "meetup-2"},
	}

	tests := []struct {
		name          string
		query         Query
		setupMock     func(repo *fakeRepo)
		expectedError error
	}{
		{
			name:  "Success: User Found and Mapped",
			query: Query{UserID: string(userID)},
			setupMock: func(repo *fakeRepo) {
				u := initialUser
				repo.users[userID] = &u
			},
			expectedError: nil,
		},
		{
			name:          "Error: User Not Found",
			query:         Query{UserID: "non-existent-id"},
			setupMock:     func(repo *fakeRepo) {},
			expectedError: shared.Error{Type: shared.ErrorTypeNotFound},
		},
		{
			name:  "Error: Database Connection Fails",
			query: Query{UserID: string(userID)},
			setupMock: func(repo *fakeRepo) {
				repo.getErr = errors.New("database connection timeout")
			},
			expectedError: errors.New("database connection timeout"),
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

				if dto.ID != string(initialUser.ID) {
					t.Errorf("expected ID %q, got %q", initialUser.ID, dto.ID)
				}
				if dto.FirstName != initialUser.FirstName {
					t.Errorf("expected FirstName %q, got %q", initialUser.FirstName, dto.FirstName)
				}
				if dto.LastName != initialUser.LastName {
					t.Errorf("expected LastName %q, got %q", initialUser.LastName, dto.LastName)
				}
				if dto.Email != initialUser.Email {
					t.Errorf("expected Email %q, got %q", initialUser.Email, dto.Email)
				}
				if dto.DisplayName != initialUser.DisplayName() {
					t.Errorf("expected DisplayName %q, got %q", initialUser.DisplayName(), dto.DisplayName)
				}
				if !reflect.DeepEqual(dto.Meetups, initialUser.Meetups) {
					t.Errorf("expected Meetups %v, got %v", initialUser.Meetups, dto.Meetups)
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
