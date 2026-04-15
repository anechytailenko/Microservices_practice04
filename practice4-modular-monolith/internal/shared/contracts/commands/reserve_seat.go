package commands

const ReserveSeatType = "meetups.seat.reserve"

type ReserveSeat struct {
	WorkflowID    string `json:"workflowId"`
	UserID        string `json:"userId"`
	MeetupID      string `json:"meetupId"`
	CorrelationID string `json:"correlationId"`
}
