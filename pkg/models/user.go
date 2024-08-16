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
