package events

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

const MeetupCreatedEventType = "events.meetups.created"

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

type MeetupCreatedEventBuilder struct {
	event MeetupCreatedEvent
}

func NewMeetupCreatedEventBuilder() *MeetupCreatedEventBuilder {
	return &MeetupCreatedEventBuilder{
		event: MeetupCreatedEvent{
			EventID:    uuid.New().String(),
			OccurredAt: time.Now().UTC(),
		},
	}
}

func (b *MeetupCreatedEventBuilder) WithCorrelationID(id string) *MeetupCreatedEventBuilder {
	b.event.CorrelationID = id
	return b
}

func (b *MeetupCreatedEventBuilder) WithMeetupData(id, ownerID, title, status string, capacity int) *MeetupCreatedEventBuilder {
	b.event.MeetupID = id
	b.event.OwnerUserID = ownerID
	b.event.Title = title
	b.event.Capacity = capacity
	b.event.Status = status
	b.event.Summary = fmt.Sprintf("Meetup '%s' created with capacity %d", title, capacity)
	return b
}

func (b *MeetupCreatedEventBuilder) Build() MeetupCreatedEvent {
	return b.event
}
