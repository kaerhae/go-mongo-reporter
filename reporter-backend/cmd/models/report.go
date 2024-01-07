package models

type Report struct {
	ID          string `bson:"_id"`
	Topic       string
	Author      string
	Description string
}
