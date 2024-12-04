package domain

import "github.com/google/uuid"

type User struct {
	GUID    uuid.UUID `gorm:"id; not null; primary_key"`
	Email   string    `gorm:"email"`
	Session Session   `gorm:"foreign:UserGUID; constraint:onUpdate:CASCADE,onDelete:CASCADE"`
}
