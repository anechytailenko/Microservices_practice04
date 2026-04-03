package create_meetup

import (
	"context"

	"github.com/anechytailenko/Microservices_practice04/internal/meetups/domain"
)

type Repository interface {
	Save(ctx context.Context, meetup *domain.Meetup) error
}

type Handler struct {
	repo Repository
}

func NewHandler(repo Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) Handle(ctx context.Context, cmd Command) (domain.MeetupID, error) {
	meetup, err := domain.NewMeetup(cmd.Title, cmd.Capacity)
	if err != nil {
		return "", err
	}

	if err := h.repo.Save(ctx, meetup); err != nil {
		return "", err
	}

	return meetup.ID, nil
}
