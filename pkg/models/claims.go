package models

import (
	"github.com/golang-jwt/jwt"
)

type Claims struct {
	Username string `json:"username"`
	AppRole  Role   `json:"appRole"`
	Token    string
	jwt.StandardClaims
}
