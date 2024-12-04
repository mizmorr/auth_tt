package domain

import "github.com/google/uuid"

type User struct {
	GUID            uuid.UUID `gorm:"id; not null; primary_key"`
	Email           string    `gorm:"email"`
	ActiveSessionID string
}

type UserRequest struct {
	Email string    `json:"email"`
	GUID  uuid.UUID `json:"guid"`
}
