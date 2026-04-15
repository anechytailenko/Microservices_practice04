package create_workflow

import (
	"context"
	"encoding/json"

	"github.com/anechytailenko/Microservices_practice04/internal/shared"
	"github.com/anechytailenko/Microservices_practice04/internal/shared/contracts/commands"
	"github.com/anechytailenko/Microservices_practice04/internal/shared/ctxutil"
	"github.com/anechytailenko/Microservices_practice04/internal/workflows/domain"
	"github.com/google/uuid"
)

type Repository interface {
	Save(ctx context.Context, workflow *domain.Workflow, eventID string, eventType string, eventPayload []byte) error
}

type Handler struct {
	repo Repository
}

func NewHandler(repo Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) Handle(ctx context.Context, cmd Command) (*WorkflowDTO, error) {
	if cmd.UserID == "" || cmd.MeetupID == "" {
		return nil, shared.NewValidationError("userId and meetupId are required")
	}

	workflow, err := domain.NewWorkflow("join-meetup")
	if err != nil {
		return nil, err
	}

	outboxPayload := commands.ReserveSeat{
		WorkflowID:    string(workflow.ID),
		UserID:        cmd.UserID,
		MeetupID:      cmd.MeetupID,
		CorrelationID: ctxutil.GetCorrelationID(ctx),
	}

	payloadBytes, err := json.Marshal(outboxPayload)
	if err != nil {
		return nil, shared.NewInternalError("failed to serialize outbox payload: %v", err)
	}

	eventID := uuid.New().String()

	eventType := commands.ReserveSeatType

	err = h.repo.Save(ctx, workflow, eventID, eventType, payloadBytes)
	if err != nil {
		return nil, err
	}

	dto := &WorkflowDTO{
		ID:        string(workflow.ID),
		Type:      workflow.Type,
		State:     string(workflow.State),
		CreatedAt: workflow.CreatedAt,
	}

	return dto, nil
}
