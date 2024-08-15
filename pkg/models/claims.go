package models

import (
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Claims struct {
	UserID      primitive.ObjectID `json:"id"`
	Username    string             `json:"username"`
	Permissions Permission         `json:"permission"`
	Token       string
	jwt.StandardClaims
}
