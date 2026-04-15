package events

const ScheduleUpdateFailedEventType = "events.users.schedule.update_failed"

type ScheduleUpdateFailed struct {
	WorkflowID    string `json:"workflowId"`
	UserID        string `json:"userId"`
	MeetupID      string `json:"meetupId"`
	Reason        string `json:"reason"`
	CorrelationID string `json:"correlationId"`
}
