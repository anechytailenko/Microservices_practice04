package domain

import (
	"time"

	"github.com/anechytailenko/Microservices_practice04/internal/shared"
)

type Meetup struct {
	ID        MeetupID
	Title     string
	Capacity  int
	Status    MeetupStatus
	CreatedAt time.Time
}

func validateDetails(title string, capacity int) error {
	if title == "" {
		return shared.NewValidationError("meetup title cannot be empty")
	}
	if capacity <= 0 {
		return shared.NewValidationError("capacity must be positive, got: %d", capacity)
	}
	return nil
}

func NewMeetup(title string, capacity int) (*Meetup, error) {
	if err := validateDetails(title, capacity); err != nil {
		return nil, err
	}

	return &Meetup{
		ID:        NewMeetupID(),
		Title:     title,
		Capacity:  capacity,
		Status:    StatusDraft,
		CreatedAt: time.Now().UTC(),
	}, nil
}

func (m *Meetup) UpdateDetails(title string, capacity int) error {
	if err := validateDetails(title, capacity); err != nil {
		return err
	}

	m.Title = title
	m.Capacity = capacity
	return nil
}

func (m *Meetup) ChangeStatus(newStatus MeetupStatus) error {
	switch m.Status {
	case StatusDraft:
		if newStatus != StatusPublished {
			return shared.NewConflictError("invalid status transition from '%s' to '%s'", m.Status, newStatus)
		}
	case StatusPublished:
		if newStatus != StatusArchived {
			return shared.NewConflictError("invalid status transition from '%s' to '%s'", m.Status, newStatus)
		}
	case StatusArchived:
		return shared.NewConflictError("invalid status transition from '%s' to '%s'", m.Status, newStatus)
	}

	m.Status = newStatus
	return nil
}
