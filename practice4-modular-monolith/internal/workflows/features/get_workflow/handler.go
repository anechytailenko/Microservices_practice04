package get_workflow

import (
	"context"

	"github.com/anechytailenko/Microservices_practice04/internal/shared"
	"github.com/anechytailenko/Microservices_practice04/internal/workflows/domain"
)

type Repository interface {
	GetByID(ctx context.Context, id domain.WorkflowID) (*domain.Workflow, error)
}

type Handler struct {
	repo Repository
}

func NewHandler(repo Repository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) Handle(ctx context.Context, q Query) (*WorkflowDTO, error) {
	if q.WorkflowID == "" {
		return nil, shared.NewValidationError("workflow_id cannot be empty")
	}

	id := domain.WorkflowID(q.WorkflowID)

	workflow, err := h.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if workflow == nil {
		return nil, shared.NewNotFoundError("workflow with id '%s' not found", q.WorkflowID)
	}

	dto := &WorkflowDTO{
		ID:        string(workflow.ID),
		Type:      workflow.Type,
		State:     string(workflow.State),
		LastError: workflow.LastError,
		CreatedAt: workflow.CreatedAt,
		UpdatedAt: workflow.UpdatedAt,
	}

	return dto, nil
}
