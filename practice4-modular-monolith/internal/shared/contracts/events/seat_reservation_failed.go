package events

<<<<<<< Updated upstream
const SeatReservationFailedEventType = "meetups.seat.reservation_failed"
=======
const SeatReservationFailedEventType = "events.meetups.seat.reservation_failed"
>>>>>>> Stashed changes

type SeatReservationFailed struct {
	WorkflowID    string `json:"workflowId"`
	UserID        string `json:"userId"`
	MeetupID      string `json:"meetupId"`
	Reason        string `json:"reason"`
	CorrelationID string `json:"correlationId"`
}
