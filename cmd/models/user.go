package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `bson:"_id"`
	Username      string             `json:"username"`
	Email         string             `json:"email"`
	Password_hash string             `json:"password_hash"`
	Created_At    string             `json:"created_at"`
	Token         string             `json:"token"`
	App_Role      string             `json:"app_role"`
}

type Role string

const (
	Admin      Role = "admin"
	Maintainer Role = "maintainer"
	Creator    Role = "creator"
	Guest      Role = "guest"
	Undefined  Role = ""
)
