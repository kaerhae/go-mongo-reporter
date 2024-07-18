package models

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	Username       string `json:"username"`
	AppRole        string `json:"appRole"`
	Token          string
	ExpirationTime time.Time
	jwt.StandardClaims
}
