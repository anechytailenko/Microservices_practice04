package get_workflow

import "time"

type WorkflowDTO struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	State     string    `json:"state"`
	LastError *string   `json:"last_error,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
