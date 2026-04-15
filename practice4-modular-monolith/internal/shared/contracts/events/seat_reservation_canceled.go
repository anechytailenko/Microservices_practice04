package events

const SeatReservationCanceledEventType = "events.meetups.seat.reservation_canceled"

type SeatReservationCanceled struct {
	WorkflowID    string `json:"workflowId"`
	MeetupID      string `json:"meetupId"`
	UserID        string `json:"userId"`
	CorrelationID string `json:"correlationId"`
}
