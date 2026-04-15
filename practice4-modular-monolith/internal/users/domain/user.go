package domain

import (
	"net/mail"
	"slices"
	"strings"

	"github.com/anechytailenko/Microservices_practice04/internal/shared"
)

type User struct {
	ID        UserID
	FirstName string
	LastName  string
	Email     string
	Meetups   []string
}

func NewUser(firstName, lastName, email string) (*User, error) {
	firstName = strings.TrimSpace(firstName)
	lastName = strings.TrimSpace(lastName)
	email = strings.TrimSpace(email)

	if firstName == "" {
		return nil, shared.NewValidationError("first name cannot be empty")
	}

	if lastName == "" {
		return nil, shared.NewValidationError("last name cannot be empty")
	}

	if email == "" {
		return nil, shared.NewValidationError("email cannot be empty")
	}

	// validation of email format pattern through library
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, shared.NewValidationError("invalid email format: %s", email)
	}

	return &User{
		ID:        NewUserID(),
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Meetups:   make([]string, 0),
	}, nil
}

func (u *User) DisplayName() string {
	return u.FirstName + " " + u.LastName
}

func (u *User) AddMeetup(meetupID string) {
	if slices.Contains(u.Meetups, meetupID) {
		return
	}
	u.Meetups = append(u.Meetups, meetupID)
}

func (u *User) RemoveMeetup(meetupID string) {
	idx := slices.Index(u.Meetups, meetupID)
	if idx == -1 {
		return
	}
	u.Meetups = slices.Delete(u.Meetups, idx, idx+1)
}
