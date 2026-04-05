package get_meetup

import "time"

type MeetupDTO struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Capacity  int       `json:"capacity"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
