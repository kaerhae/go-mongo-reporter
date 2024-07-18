package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Report struct {
	ID          primitive.ObjectID `bson:"_id"`
	Topic       string             `json:"topic"`
	Author      string             `json:"author"`
	Description string             `json:"description"`
	UserID      string             `bson:"userID"`
}
