package commands

const SendNotificationType = "commands.notifications.send"

type SendNotification struct {
	WorkflowID    string `json:"workflowId"`
	UserID        string `json:"userId"`
	MeetupID      string `json:"meetupId"`
	Status        string `json:"status"`
	Message       string `json:"message"`
	CorrelationID string `json:"correlationId"`
}
