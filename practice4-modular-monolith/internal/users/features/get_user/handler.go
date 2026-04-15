package get_user

import (
	"context"

	"github.com/anechytailenko/Microservices_practice04/internal/shared"
	"github.com/anechytailenko/Microservices_practice04/internal/users/domain"
)

type Repository interface {
	GetByID(ctx context.Context, id domain.UserID) (*domain.User, error)
}

type Handler struct {
	repo Repository
}

func NewHandler(repo Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) Handle(ctx context.Context, q Query) (*UserDTO, error) {
	id := domain.UserID(q.UserID)

	user, err := h.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, shared.NewNotFoundError("user with id '%s' not found", q.UserID)
	}

	dto := &UserDTO{
		ID:          string(user.ID),
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		DisplayName: user.DisplayName(),
		Meetups:     user.Meetups,
	}

	return dto, nil
}
