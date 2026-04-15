package events

const ScheduleUpdatedEventType = "events.users.schedule.updated"

type ScheduleUpdated struct {
	WorkflowID    string `json:"workflowId"`
	UserID        string `json:"userId"`
	MeetupID      string `json:"meetupId"`
	CorrelationID string `json:"correlationId"`
}
