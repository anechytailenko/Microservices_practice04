package domain

import "github.com/google/uuid"

type WorkflowID string

func NewWorkflowID() WorkflowID {
	return WorkflowID(uuid.New().String())
}
