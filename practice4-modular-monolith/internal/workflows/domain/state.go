package domain

type WorkflowState string

const (
	StateInitializing       WorkflowState = "Initializing"
	StateSeatReserved       WorkflowState = "SeatReserved"
<<<<<<< Updated upstream
	StateScheduleUpdated    WorkflowState = "ScheduleUpdated"
=======
>>>>>>> Stashed changes
	StateCancelingSeat      WorkflowState = "CancelingSeat"
	StateCompleted          WorkflowState = "Completed"
	StateFailed             WorkflowState = "Failed"
	StateManualIntervention WorkflowState = "ManualInterventionRequired"
)
