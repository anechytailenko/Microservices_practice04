package events

const SeatReservationFailedEventType = "meetups.seat.reservation_failed"

type SeatReservationFailed struct {
	WorkflowID    string `json:"workflowId"`
	UserID        string `json:"userId"`
	MeetupID      string `json:"meetupId"`
	Reason        string `json:"reason"`
	CorrelationID string `json:"correlationId"`
}
