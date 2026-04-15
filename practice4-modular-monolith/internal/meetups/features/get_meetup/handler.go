package get_meetup

import (
	"context"

	"github.com/anechytailenko/Microservices_practice04/internal/meetups/domain"
	"github.com/anechytailenko/Microservices_practice04/internal/shared"
)

type Repository interface {
	GetByID(ctx context.Context, id domain.MeetupID) (*domain.Meetup, error)
}

type Handler struct {
	repo Repository
}

func NewHandler(repo Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) Handle(ctx context.Context, q Query) (*MeetupDTO, error) {
	id := domain.MeetupID(q.MeetupID)

	meetup, err := h.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if meetup == nil {
		return nil, shared.NewNotFoundError("meetup with id '%s' not found", q.MeetupID)
	}

	dto := &MeetupDTO{
		ID:          string(meetup.ID),
		Title:       meetup.Title,
		Capacity:    meetup.Capacity,
		OwnerUserID: meetup.OwnerUserID,
		Status:      string(meetup.Status),
		Guests:      meetup.Guests,
		CreatedAt:   meetup.CreatedAt,
	}

	return dto, nil
}
