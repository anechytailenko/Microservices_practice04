package events

<<<<<<< Updated upstream
const SeatReservationCanceledEventType = "meetups.seat.reservation_canceled"
=======
const SeatReservationCanceledEventType = "events.meetups.seat.reservation_canceled"
>>>>>>> Stashed changes

type SeatReservationCanceled struct {
	WorkflowID    string `json:"workflowId"`
	MeetupID      string `json:"meetupId"`
	UserID        string `json:"userId"`
	CorrelationID string `json:"correlationId"`
}
