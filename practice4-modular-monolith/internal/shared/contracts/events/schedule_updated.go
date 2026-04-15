package events

<<<<<<< Updated upstream
const ScheduleUpdatedEventType = "users.schedule.updated"
=======
const ScheduleUpdatedEventType = "events.users.schedule.updated"
>>>>>>> Stashed changes

type ScheduleUpdated struct {
	WorkflowID    string `json:"workflowId"`
	UserID        string `json:"userId"`
	MeetupID      string `json:"meetupId"`
	CorrelationID string `json:"correlationId"`
}
