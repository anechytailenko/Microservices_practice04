package create_user

import (
	"context"

	"github.com/anechytailenko/Microservices_practice04/internal/users/domain"
)

type Repository interface {
	Save(ctx context.Context, user *domain.User) error
}

type Handler struct {
	repo Repository
}

func NewHandler(repo Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) Handle(ctx context.Context, cmd Command) (domain.UserID, error) {
	user, err := domain.NewUser(cmd.FirstName, cmd.LastName, cmd.Email)
	if err != nil {
		return "", err
	}

	if err := h.repo.Save(ctx, user); err != nil {
		return "", err
	}

	return user.ID, nil
}
