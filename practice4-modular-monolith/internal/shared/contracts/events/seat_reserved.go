package events

<<<<<<< Updated upstream
const SeatReservedEventType = "meetups.seat.reserved"
=======
const SeatReservedEventType = "events.meetups.seat.reserved"
>>>>>>> Stashed changes

type SeatReserved struct {
	WorkflowID    string `json:"workflowId"`
	MeetupID      string `json:"meetupId"`
	UserID        string `json:"userId"`
	CorrelationID string `json:"correlationId"`
}
