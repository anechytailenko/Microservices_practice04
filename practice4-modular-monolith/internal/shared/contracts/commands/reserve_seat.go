package commands

<<<<<<< Updated upstream
const ReserveSeatType = "meetups.seat.reserve"
=======
const ReserveSeatType = "commands.meetups.seat.reserve"
>>>>>>> Stashed changes

type ReserveSeat struct {
	WorkflowID    string `json:"workflowId"`
	UserID        string `json:"userId"`
	MeetupID      string `json:"meetupId"`
	CorrelationID string `json:"correlationId"`
}
