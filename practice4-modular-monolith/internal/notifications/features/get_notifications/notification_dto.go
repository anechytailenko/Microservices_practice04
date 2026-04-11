package get_notifications

import "time"

type NotificationDTO struct {
	MeetupID  string    `json:"meetup_id"`
	Title     string    `json:"title"`
	Summary   string    `json:"summary"`
	Capacity  int       `json:"capacity"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
