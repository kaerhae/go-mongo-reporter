package services

import (
	"errors"
	"fmt"
	"log"
	"main/cmd/db"
	"main/cmd/models"
	"main/configs"
	"time"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

var secretKey = []byte(configs.GetSecret())

func CheckExistingUser(username string) (models.User, error) {

	client, ctx, cancel, err := db.MongoConnect(configs.GetMongoURI())
	if err != nil {
		return models.User{}, err
	}

	collection := client.Database(configs.GetDBName()).Collection("users")
	filter := bson.D{{Key: "username", Value: username}}
	var result models.User
	err = collection.FindOne(ctx, filter).Decode(&result)

	if err != nil {
		return models.User{}, err
	}

	defer db.MongoClose(client, ctx, cancel)

	return result, nil

}

func CreateUser(user models.User) (models.User, error) {

	client, ctx, cancel, err := db.MongoConnect(configs.GetMongoURI())
	if err != nil {
		return models.User{}, err
	}

	collection := client.Database(configs.GetDBName()).Collection("users")

	addUser, err := collection.InsertOne(ctx, user)
	if err != nil {
		return models.User{}, err
	}

	fmt.Println("Added new user: ", addUser.InsertedID)

	defer db.MongoClose(client, ctx, cancel)

	return user, nil
}

func HashPwd(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.MinCost,
	)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return string(hash), nil
}

func CheckPassword(hashedPassword string, plainPassword string) bool {
	bytePlain := []byte(plainPassword)
	byteHashed := []byte(hashedPassword)

	err := bcrypt.CompareHashAndPassword(byteHashed, bytePlain)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func CreateToken(username string) (*models.Claims, error) {
	expirationTime := time.Now().Add(5 * time.Minute)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Minute * 30).Unix(),
		})

	tokenstring, err := token.SignedString(secretKey)
	claims := &models.Claims{
		Username:       username,
		Token:          tokenstring,
		ExpirationTime: expirationTime,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	if err != nil {
		return nil, err
	}

	return claims, nil
}

func DetermineRole(role string) (models.Role, error) {
	fmt.Println(role)
	switch role {
	case "admin":
		return models.Admin, nil
	case "maintainer":
		return models.Maintainer, nil
	case "creator":
		return models.Creator, nil
	case "guest":
		return models.Guest, nil
	default:
		return models.Undefined, errors.New("Role undefined: " + role)
	}
}
