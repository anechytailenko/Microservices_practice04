package get_meetup

import "time"

type MeetupDTO struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Capacity    int       `json:"capacity"`
	OwnerUserID string    `json:"owner_user_id"`
	Guests      []string  `json:"guests"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}
