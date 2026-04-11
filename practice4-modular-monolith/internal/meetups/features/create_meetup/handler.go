package create_meetup

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/anechytailenko/Microservices_practice04/internal/meetups/domain"
	"github.com/anechytailenko/Microservices_practice04/internal/shared/ctxutil"
	"github.com/anechytailenko/Microservices_practice04/internal/shared/events"
)

type Repository interface {
	Save(ctx context.Context, meetup *domain.Meetup, eventID string, eventType string, eventPayload []byte) error
}

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

	evt := events.NewMeetupCreatedEventBuilder().
		WithCorrelationID(ctxutil.GetCorrelationID(ctx)).
		WithMeetupData(
			string(meetup.ID),
			meetup.OwnerUserID,
			meetup.Title,
			string(meetup.Status),
			meetup.Capacity,
		).
		Build()

	payloadBytes, err := json.Marshal(evt)
	if err != nil {
		return "", fmt.Errorf("failed to marshal outbox event: %w", err)
	}

	if err := h.repo.Save(ctx, meetup, evt.EventID, "meetup.created", payloadBytes); err != nil {
		return "", err
	}

	return meetup.ID, nil
}
