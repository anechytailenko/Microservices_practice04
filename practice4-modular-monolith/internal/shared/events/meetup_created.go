package events

import "time"

type MeetupCreatedEvent struct {
	EventID       string    `json:"eventId"`
	OccurredAt    time.Time `json:"occurredAt"`
	CorrelationID string    `json:"correlationId"`
	MeetupID      string    `json:"coreItemId"`
	OwnerUserID   string    `json:"ownerUserId"`
	Summary       string    `json:"summary"`
	Title         string    `json:"title"`
	Capacity      int       `json:"capacity"`
	Status        string    `json:"status"`
}
