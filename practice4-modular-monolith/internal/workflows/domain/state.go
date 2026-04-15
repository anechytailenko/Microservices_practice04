package domain

type WorkflowState string

const (
	StateInitializing       WorkflowState = "Initializing"
	StateSeatReserved       WorkflowState = "SeatReserved"
	StateScheduleUpdated    WorkflowState = "ScheduleUpdated"
	StateCancelingSeat      WorkflowState = "CancelingSeat"
	StateCompleted          WorkflowState = "Completed"
	StateFailed             WorkflowState = "Failed"
	StateManualIntervention WorkflowState = "ManualInterventionRequired"
)
