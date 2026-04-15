package commands

<<<<<<< Updated upstream
const UpdateScheduleType = "users.schedule.update"
=======
const UpdateScheduleType = "commands.users.schedule.update"
>>>>>>> Stashed changes

type UpdateSchedule struct {
	WorkflowID    string `json:"workflowId"`
	UserID        string `json:"userId"`
	MeetupID      string `json:"meetupId"`
	CorrelationID string `json:"correlationId"`
}
