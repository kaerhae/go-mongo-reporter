package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
<<<<<<< HEAD
	ID            primitive.ObjectID `bson:"_id"`
	Username      string             `json:"username"`
	Email         string             `json:"email"`
	Password_hash string             `json:"password_hash"`
	Created_At    string             `json:"created_at"`
	Token         string             `json:"token"`
	App_Role      string             `json:"app_role"`
=======
	ID           primitive.ObjectID `bson:"_id"`
	Username     string             `bson:"username"`
	Email        string             `bson:"email"`
	PasswordHash string             `bson:"password_hash"`
	CreatedAt    string             `bson:"created_at"`
	Token        string             `bson:"token"`
	AppRole      string             `bson:"app_role"`
>>>>>>> master
}

type Role string

const (
	Admin      Role = "admin"
	Maintainer Role = "maintainer"
	Creator    Role = "creator"
	Guest      Role = "guest"
	Undefined  Role = ""
)
