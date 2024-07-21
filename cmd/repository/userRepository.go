package repository

import (
	"context"
	"main/cmd/middleware"
	"main/cmd/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	Create(user *models.User) (string, error)
	GetSingleUser(username string) (models.User, error)
}

type userRepository struct {
	Client *mongo.Database
	Logger middleware.Logger
}

func NewUserRepository(client *mongo.Database, logger middleware.Logger) UserRepository {
	return &userRepository{
		Client: client,
		Logger: logger,
	}
}

func (r *userRepository) Create(user *models.User) (string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	addUser, err := r.Client.Collection("users").InsertOne(ctx, &user)
	if err != nil {
		return "", err
	}

	t := addUser.InsertedID.(primitive.ObjectID).Hex()
	return t, nil
}

// GetSingleUser implements UserRepository.
func (r *userRepository) GetSingleUser(username string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	collection := r.Client.Collection("users")
	filter := bson.D{{Key: "username", Value: username}}
	var result models.User
	err := collection.FindOne(ctx, filter).Decode(&result)

	if err != nil {
		return models.User{}, err
	}

	return result, nil
}
