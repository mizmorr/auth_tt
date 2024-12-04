package domain

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type AccessClaims struct {
	GUID uuid.UUID `json:"guid"`
	IP   string    `json:"ip"`
	jwt.StandardClaims
}
