package commands

<<<<<<< Updated upstream
const CancelSeatReservationType = "meetups.seat.cancel"
=======
const CancelSeatReservationType = "commands.meetups.seat.cancel"
>>>>>>> Stashed changes

type CancelSeatReservation struct {
	WorkflowID    string `json:"workflowId"`
	UserID        string `json:"userId"`
	MeetupID      string `json:"meetupId"`
	CorrelationID string `json:"correlationId"`
}
