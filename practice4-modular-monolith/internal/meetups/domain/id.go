package domain

import "github.com/google/uuid"

type MeetupID string

func NewMeetupID() MeetupID {
	return MeetupID(uuid.New().String())
}
