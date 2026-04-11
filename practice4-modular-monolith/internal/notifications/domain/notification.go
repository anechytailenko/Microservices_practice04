package domain

import "time"

type Notification struct {
	ID          string
	OwnerUserID string
	MeetupID    string
	Title       string
	Summary     string
	Capacity    int
	Status      string
	CreatedAt   time.Time
}
