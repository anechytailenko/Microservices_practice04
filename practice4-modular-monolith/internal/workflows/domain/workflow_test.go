package domain

import (
	"testing"
)

func TestNewWorkflow(t *testing.T) {
	tests := []struct {
		name         string
		workflowType string
		expectError  bool
	}{
		{"Valid Join Meetup Workflow", "join-meetup", false},
		{"Empty Type Workflow", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, err := NewWorkflow(tt.workflowType)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected validation error, got none")
				}
				if w != nil {
					t.Errorf("expected workflow to be nil on error, got object")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if w.ID == "" {
					t.Errorf("expected workflow ID to be generated")
				}
				if w.State != StateInitializing {
					t.Errorf("expected initial state %q, got %q", StateInitializing, w.State)
				}
				if w.LastError != nil {
					t.Errorf("expected last error to be nil, got %v", *w.LastError)
				}
			}
		})
	}
}

func TestWorkflow_ChangeState(t *testing.T) {
	tests := []struct {
		name        string
		startState  WorkflowState
		newState    WorkflowState
		expectError bool
	}{
		{"Initializing -> SeatReserved", StateInitializing, StateSeatReserved, false},
		{"SeatReserved -> Completed", StateSeatReserved, StateCompleted, false},
		{"SeatReserved -> CancelingSeat", StateSeatReserved, StateCancelingSeat, false},
		{"CancelingSeat -> Failed", StateCancelingSeat, StateFailed, false},

		{"Initializing -> Completed (Skip steps)", StateInitializing, StateCompleted, true},
		{"Completed -> Failed (Modify terminal)", StateCompleted, StateFailed, true},
		{"Failed -> Initializing (Restart from terminal)", StateFailed, StateInitializing, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Workflow{
				State: tt.startState,
			}

			err := w.ChangeState(tt.newState)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error for transition %s -> %s, got none", tt.startState, tt.newState)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for transition %s -> %s: %v", tt.startState, tt.newState, err)
				}
				if w.State != tt.newState {
					t.Errorf("expected state to change to %s, got %s", tt.newState, w.State)
				}
			}
		})
	}
}
