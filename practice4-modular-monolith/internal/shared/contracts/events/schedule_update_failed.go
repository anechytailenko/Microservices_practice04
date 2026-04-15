package events

<<<<<<< Updated upstream
const ScheduleUpdateFailedEventType = "users.schedule.update_failed"
=======
const ScheduleUpdateFailedEventType = "events.users.schedule.update_failed"
>>>>>>> Stashed changes

type ScheduleUpdateFailed struct {
	WorkflowID    string `json:"workflowId"`
	UserID        string `json:"userId"`
	MeetupID      string `json:"meetupId"`
	Reason        string `json:"reason"`
	CorrelationID string `json:"correlationId"`
}
