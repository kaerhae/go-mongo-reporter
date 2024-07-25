package repository

import (
	"context"
	"main/pkg/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository interface {
	Create(user *models.User) (string, error)
	Get() ([]models.User, error)
	GetSingleUser(username string) (models.User, error)
	UpdateSingleUser(user models.User) error
	DeleteSingleUser(ID string) (int64, error)
}

type userRepository struct {
	Client *mongo.Database
}

// Get implements UserRepository.
func (r *userRepository) Get() ([]models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	collection := r.Client.Collection("users")
	var users []models.User

	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	if err = cur.All(context.TODO(), &users); err != nil {
		panic(err)
	}

	return users, nil
}

// UpdateSingleUser implements UserRepository.
func (r *userRepository) UpdateSingleUser(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	collection := r.Client.Collection("users")

	_, err := collection.UpdateOne(ctx, bson.D{{
		Key:   "_id",
		Value: user.ID,
	}}, bson.M{"$set": user})

	if err != nil {
		return err
	}

	return nil
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

// DeleteSingleUser implements UserRepository
func (r *userRepository) DeleteSingleUser(id string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	collection := r.Client.Collection("users")

	deleteCount, err := collection.DeleteOne(ctx, bson.D{{
		Key:   "_id",
		Value: id,
	}})
	if err != nil {
		return 0, err
	}

	return deleteCount.DeletedCount, nil
}

func NewUserRepository(client *mongo.Database) UserRepository {
	return &userRepository{
		Client: client,
	}
}
