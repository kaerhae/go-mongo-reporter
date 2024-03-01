package models

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	Username       string `json:"username"`
	App_Role       string `json:"app_role"`
	Token          string
	ExpirationTime time.Time
	jwt.StandardClaims
}
