package models

import (
	"errors"
	"fmt"

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

func DetermineRole(role string) (Role, error) {
	fmt.Println(role)
	switch role {
	case "admin":
		return Admin, nil
	case "maintainer":
		return Maintainer, nil
	case "creator":
		return Creator, nil
	case "guest":
		return Guest, nil
	default:
		return Undefined, errors.New("Role undefined: " + role)
	}
}
