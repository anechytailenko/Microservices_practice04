package change_status

import (
	"context"

	"github.com/anechytailenko/Microservices_practice04/internal/meetups/domain"
	"github.com/anechytailenko/Microservices_practice04/internal/shared"
)

type Repository interface {
	GetByID(ctx context.Context, id domain.MeetupID) (*domain.Meetup, error)
	Update(ctx context.Context, meetup *domain.Meetup) error
}

type Handler struct {
	repo Repository
}

func NewHandler(repo Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) Handle(ctx context.Context, cmd Command) error {
	id := domain.MeetupID(cmd.MeetupID)
	newStatus := domain.MeetupStatus(cmd.Status)

	meetup, err := h.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if meetup == nil {
		return shared.NewNotFoundError("meetup with id '%s' not found", cmd.MeetupID)
	}

	if err := meetup.ChangeStatus(newStatus); err != nil {
		return err
	}

	if err := h.repo.Update(ctx, meetup); err != nil {
		return err
	}

	return nil
}
