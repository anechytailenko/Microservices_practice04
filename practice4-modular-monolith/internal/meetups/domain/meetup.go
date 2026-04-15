package domain

import (
	"slices"
	"time"

	"github.com/anechytailenko/Microservices_practice04/internal/shared"
)

type Meetup struct {
	ID          MeetupID
	Title       string
	Capacity    int
	OwnerUserID string
	Guests      []string
	Status      MeetupStatus
	CreatedAt   time.Time
}

func NewMeetup(title string, capacity int, ownerUserID string) (*Meetup, error) {
	if title == "" {
		return nil, shared.NewValidationError("meetup title cannot be empty")
	}
	if capacity <= 0 {
		return nil, shared.NewValidationError("capacity must be positive, got: %d", capacity)
	}
	if ownerUserID == "" {
		return nil, shared.NewValidationError("owner user id cannot be empty")
	}

	return &Meetup{
		ID:          NewMeetupID(),
		Title:       title,
		Capacity:    capacity,
		OwnerUserID: ownerUserID,
		Guests:      make([]string, 0),
		Status:      StatusDraft,
		CreatedAt:   time.Now().UTC(),
	}, nil
}

func (m *Meetup) ChangeStatus(newStatus MeetupStatus) error {
	switch m.Status {
	case StatusDraft:
		if newStatus != StatusPublished && newStatus != StatusCanceled {
			return shared.NewConflictError("invalid status transition from '%s' to '%s'", m.Status, newStatus)
		}
	case StatusPublished:
		if newStatus != StatusArchived && newStatus != StatusCanceled {
			return shared.NewConflictError("invalid status transition from '%s' to '%s'", m.Status, newStatus)
		}
	case StatusArchived, StatusCanceled:
		return shared.NewConflictError("invalid status transition from '%s' to '%s'", m.Status, newStatus)
	}

	m.Status = newStatus
	return nil
}

func (m *Meetup) AddGuest(userID string) error {
	// domain idempotency
	if slices.Contains(m.Guests, userID) {
		return nil
	}

	if len(m.Guests) >= m.Capacity {
		return shared.NewConflictError("meetup is fully booked")
	}

	m.Guests = append(m.Guests, userID)
	return nil
}

func (m *Meetup) RemoveGuest(userID string) {
	// domain idempotency
	idx := slices.Index(m.Guests, userID)
	if idx == -1 {
		return
	}

	m.Guests = slices.Delete(m.Guests, idx, idx+1)
}
