package events

const SeatReservedEventType = "events.meetups.seat.reserved"

type SeatReserved struct {
	WorkflowID    string `json:"workflowId"`
	MeetupID      string `json:"meetupId"`
	UserID        string `json:"userId"`
	CorrelationID string `json:"correlationId"`
}
