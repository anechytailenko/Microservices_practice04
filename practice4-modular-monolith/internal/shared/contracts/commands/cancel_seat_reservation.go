package commands

const CancelSeatReservationType = "meetups.seat.cancel"

type CancelSeatReservation struct {
	WorkflowID    string `json:"workflowId"`
	UserID        string `json:"userId"`
	MeetupID      string `json:"meetupId"`
	CorrelationID string `json:"correlationId"`
}
