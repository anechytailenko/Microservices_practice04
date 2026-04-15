package get_meetup

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/anechytailenko/Microservices_practice04/internal/meetups/domain"
	"github.com/anechytailenko/Microservices_practice04/internal/shared"
)

type fakeRepo struct {
	meetups map[domain.MeetupID]*domain.Meetup
	getErr  error
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

func TestHandler_Handle(t *testing.T) {
	meetupID := domain.MeetupID("query-123")
	now := time.Now().UTC()

	initialMeetup := domain.Meetup{
		ID:          meetupID,
		Title:       "Read API Architecture",
		Capacity:    200,
		OwnerUserID: "organizer-999",
		Guests:      []string{"user-1", "user-2"},
		Status:      domain.StatusPublished,
		CreatedAt:   now,
	}

	tests := []struct {
		name          string
		query         Query
		setupMock     func(repo *fakeRepo)
		expectedError error
	}{
		{
			name:  "Success: Meetup Found and Mapped",
			query: Query{MeetupID: string(meetupID)},
			setupMock: func(repo *fakeRepo) {
				m := initialMeetup
				repo.meetups[meetupID] = &m
			},
			expectedError: nil,
		},
		{
			name:          "Error: Meetup Not Found",
			query:         Query{MeetupID: "non-existent-id"},
			setupMock:     func(repo *fakeRepo) {},
			expectedError: shared.Error{Type: shared.ErrorTypeNotFound},
		},
		{
			name:  "Error: Database Connection Fails",
			query: Query{MeetupID: string(meetupID)},
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

				if dto.ID != string(initialMeetup.ID) {
					t.Errorf("expected ID %q, got %q", initialMeetup.ID, dto.ID)
				}
				if dto.Title != initialMeetup.Title {
					t.Errorf("expected Title %q, got %q", initialMeetup.Title, dto.Title)
				}
				if dto.Capacity != initialMeetup.Capacity {
					t.Errorf("expected Capacity %d, got %d", initialMeetup.Capacity, dto.Capacity)
				}
				if dto.OwnerUserID != initialMeetup.OwnerUserID {
					t.Errorf("expected OwnerUserID %q, got %q", initialMeetup.OwnerUserID, dto.OwnerUserID)
				}
				if !reflect.DeepEqual(dto.Guests, initialMeetup.Guests) {
					t.Errorf("expected Guests %v, got %v", initialMeetup.Guests, dto.Guests)
				}
				if dto.Status != string(initialMeetup.Status) {
					t.Errorf("expected Status %q, got %q", initialMeetup.Status, dto.Status)
				}
				if !dto.CreatedAt.Equal(initialMeetup.CreatedAt) {
					t.Errorf("expected CreatedAt %v, got %v", initialMeetup.CreatedAt, dto.CreatedAt)
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
