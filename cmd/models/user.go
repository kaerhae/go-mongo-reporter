package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `bson:"_id"`
	Username      string
	Email         string
	Password_hash string
	Created_At    string
	Token         string
	App_Role      string
}

type Role string

const (
	Admin      Role = "admin"
	Maintainer Role = "maintainer"
	Creator    Role = "creator"
	Guest      Role = "guest"
	Undefined  Role = ""
)
