package services

import (
	"fmt"
	"log"
	"main/cmd/db"
	"main/cmd/models"
	"main/configs"

	"go.mongodb.org/mongo-driver/bson"
)

func CheckExistingUser(username string) bool {

	client, ctx, cancel, err := db.MongoConnect(configs.GetMongoURI())
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database(configs.GetDBName()).Collection("users")
	filter := bson.D{{Key: "username", Value: username}}
	var result models.User
	err = collection.FindOne(ctx, filter).Decode(&result)

	if err != nil {
		return false
	}

	defer db.MongoClose(client, ctx, cancel)

	return true

}

func CreateUser(user models.User) models.User {

	client, ctx, cancel, err := db.MongoConnect(configs.GetMongoURI())
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database(configs.GetDBName()).Collection("users")

	addUser, err := collection.InsertOne(ctx, user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Added new user: ", addUser.InsertedID)

	defer db.MongoClose(client, ctx, cancel)

	return user
}
