package service

import (
	"cmd/reporter-backend/cmd/databaseContext"
	"cmd/reporter-backend/cmd/models"
	"cmd/reporter-backend/configs"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

func CheckExistingUser(username string) bool {

	client, ctx, cancel, err := databaseContext.MongoConnect(configs.GetMongoURI())
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

	defer databaseContext.MongoClose(client, ctx, cancel)

	return true

}

func CreateUser(user models.User) models.User {

	client, ctx, cancel, err := databaseContext.MongoConnect(configs.GetMongoURI())
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database(configs.GetDBName()).Collection("users")

	addUser, err := collection.InsertOne(ctx, user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Added new user: ", addUser.InsertedID)

	defer databaseContext.MongoClose(client, ctx, cancel)

	return user
}
