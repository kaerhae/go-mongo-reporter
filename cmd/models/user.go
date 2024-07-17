package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id"`
	Username     string             `bson:"username"`
	Email        string             `bson:"email"`
	PasswordHash string             `bson:"password_hash"`
	CreatedAt    string             `bson:"created_at"`
	Token        string             `bson:"token"`
	AppRole      string             `bson:"app_role"`
}

type Role string

const (
	Admin      Role = "admin"
	Maintainer Role = "maintainer"
	Creator    Role = "creator"
	Guest      Role = "guest"
	Undefined  Role = ""
)
