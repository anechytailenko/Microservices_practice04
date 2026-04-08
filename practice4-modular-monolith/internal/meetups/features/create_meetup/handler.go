package create_meetup

import (
	"context"

	"github.com/anechytailenko/Microservices_practice04/internal/meetups/domain"
)

type Repository interface {
	Save(ctx context.Context, meetup *domain.Meetup) error
}

// interface that is implemented by http-client
type UserValidator interface {
	ValidateUserExists(ctx context.Context, userID string) error
}

type Handler struct {
	repo          Repository
	userValidator UserValidator
}

func NewHandler(repo Repository, userValidator UserValidator) *Handler {
	return &Handler{
		repo:          repo,
		userValidator: userValidator,
	}
}

func (h *Handler) Handle(ctx context.Context, cmd Command) (domain.MeetupID, error) {
	if err := h.userValidator.ValidateUserExists(ctx, cmd.OwnerUserID); err != nil {
		return "", err
	}

	meetup, err := domain.NewMeetup(cmd.Title, cmd.Capacity, cmd.OwnerUserID)
	if err != nil {
		return "", err
	}

	if err := h.repo.Save(ctx, meetup); err != nil {
		return "", err
	}

	return meetup.ID, nil
}
