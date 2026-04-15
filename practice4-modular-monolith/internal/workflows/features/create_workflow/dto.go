package create_workflow

import "time"

type WorkflowDTO struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	State     string    `json:"state"`
	CreatedAt time.Time `json:"created_at"`
}
