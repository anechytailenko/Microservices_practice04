package domain

type MeetupStatus string

const (
	StatusDraft     MeetupStatus = "Draft"
	StatusPublished MeetupStatus = "Published"
	StatusArchived  MeetupStatus = "Archived"
	StatusCanceled  MeetupStatus = "Canceled"
)
