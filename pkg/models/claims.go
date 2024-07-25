package models

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	Username       string `json:"username"`
	AppRole        Role   `json:"appRole"`
	Token          string
	ExpirationTime time.Time
	jwt.StandardClaims
}
