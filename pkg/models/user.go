package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID   `bson:"_id"`
	Username     string               `bson:"username"`
	Email        string               `bson:"email"`
	PasswordHash string               `bson:"password_hash"`
	CreatedAt    string               `bson:"created_at"`
	Token        string               `bson:"token"`
	AppRole      Role                 `bson:"app_role"`
	Reports      []primitive.ObjectID `bson:"reports"`
}

type Role string

const (
	Admin      Role = "admin"
	Maintainer Role = "maintainer"
	Creator    Role = "creator"
	Guest      Role = "guest"
	Undefined  Role = ""
)

type CreateUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	AppRole  string `json:"appRole"`
}

type LoginUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
