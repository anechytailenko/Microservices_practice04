package domain

import "github.com/google/uuid"

type UserID string

func NewUserID() UserID {
	return UserID(uuid.New().String())
}
