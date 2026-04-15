package domain

import (
	"time"

	"github.com/anechytailenko/Microservices_practice04/internal/shared"
)

type Workflow struct {
	ID        WorkflowID
	Type      string
	State     WorkflowState
	LastError *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewWorkflow(workflowType string) (*Workflow, error) {
	if workflowType == "" {
		return nil, shared.NewValidationError("workflow type cannot be empty")
	}

	now := time.Now().UTC()

	return &Workflow{
		ID:        NewWorkflowID(),
		Type:      workflowType,
		State:     StateInitializing,
		LastError: nil,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (w *Workflow) ChangeState(newState WorkflowState) error {
	if !w.isValidTransition(newState) {
		return shared.NewConflictError("invalid workflow state transition from '%s' to '%s'", w.State, newState)
	}

	w.State = newState
	w.UpdatedAt = time.Now().UTC()
	return nil
}

func (w *Workflow) MarkAsFailed(failedState WorkflowState, reason string) error {
	if err := w.ChangeState(failedState); err != nil {
		return err
	}

	w.LastError = &reason
	return nil
}

func (w *Workflow) isValidTransition(newState WorkflowState) bool {
	switch w.State {
	case StateInitializing:
		return newState == StateSeatReserved || newState == StateFailed

	case StateSeatReserved:
<<<<<<< Updated upstream
		return newState == StateScheduleUpdated || newState == StateCancelingSeat

	case StateScheduleUpdated:
		return newState == StateCompleted
=======
		return newState == StateCompleted || newState == StateCancelingSeat
>>>>>>> Stashed changes

	case StateCancelingSeat:
		return newState == StateFailed || newState == StateManualIntervention

	case StateCompleted, StateFailed, StateManualIntervention:
		return false

	default:
		return false
	}
}
