package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID   `json:"id" bson:"_id"`
	Username     string               `json:"username" bson:"username"`
	Email        string               `json:"email" bson:"email"`
	PasswordHash string               `json:"-" bson:"password_hash"`
	CreatedAt    string               `json:"createdAt" bson:"created_at"`
	Permission   Permission           `json:"permission" bson:"permission"`
	Reports      []primitive.ObjectID `json:"reports" bson:"reports"`
}

type CreateUser struct {
	Username   string     `json:"username"`
	Email      string     `json:"email"`
	Password   string     `json:"password"`
	Permission Permission `json:"permission"`
}

type CreateGuestUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
