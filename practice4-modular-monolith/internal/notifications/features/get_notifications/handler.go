package get_notifications

import (
	"context"

	"github.com/anechytailenko/Microservices_practice04/internal/notifications/domain"
)

type Repository interface {
	GetByOwnerID(ctx context.Context, ownerUserID string) ([]domain.Notification, error)
}

type Handler struct {
	repo Repository
}

func NewHandler(repo Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) Handle(ctx context.Context, q Query) ([]NotificationDTO, error) {
	domainNotifications, err := h.repo.GetByOwnerID(ctx, q.OwnerUserID)
	if err != nil {
		return nil, err
	}

	if len(domainNotifications) == 0 {
		return []NotificationDTO{}, nil
	}

	var dtos []NotificationDTO

	for _, n := range domainNotifications {
		dtos = append(dtos, NotificationDTO{
			MeetupID:  n.MeetupID,
			Title:     n.Title,
			Summary:   n.Summary,
			Capacity:  n.Capacity,
			Status:    n.Status,
			CreatedAt: n.CreatedAt,
		})
	}

	return dtos, nil
}
