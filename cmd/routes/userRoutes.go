package routes

import (
	"fmt"
	"main/cmd/middleware"
	"main/cmd/models"
	"main/cmd/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRouter interface {
	PostNewUser(c *gin.Context)
	LoginUser(c *gin.Context)
}

type userRouter struct {
	Service services.UserService
	Logger  middleware.Logger
}

var validate = validator.New()

func NewUserHandler(service services.UserService, logger middleware.Logger) UserRouter {
	return &userRouter{
		Service: service,
		Logger:  logger,
	}
}

// POST /users
func (u *userRouter) PostNewUser(c *gin.Context) {
	var body models.User
	/* Bindataan request body muuttujaan body */
	if err := c.BindJSON(&body); err != nil {
		u.Logger.LogError(
			fmt.Sprintf("Error happened while binding JSON: %v", err),
		)
		c.IndentedJSON(http.StatusBadRequest, "Error in handling request")
		return
	}

	_, err := u.Service.CheckExistingUser(body.Username)

	/* Tarkistetaan löytyykö käyttäjää ennestään */
	if err == nil {

		c.IndentedJSON(http.StatusBadRequest, "Username already exists")
		return
	}

	if validationErr := validate.Struct(&body); validationErr != nil {
		u.Logger.LogError(
			fmt.Sprintf("Error happened while validating body: %v", err),
		)
		c.IndentedJSON(http.StatusBadRequest, "Malformatted body")
		return
	}

	role, err := u.Service.DetermineRole(string(body.AppRole))
	if err != nil {
		u.Logger.LogError(
			fmt.Sprintf("Malformatted appRole: %v", err),
		)
		c.IndentedJSON(http.StatusBadRequest, "Malformatted role")
		return
	}

	hash, err := u.Service.HashPwd(body.PasswordHash)
	if err != nil {
		u.Logger.LogError(
			fmt.Sprintf("Error happened while hashing password: %v", err),
		)
		c.IndentedJSON(500, "Server failed")
		return
	}

	newUser := models.User{
		ID:           primitive.NewObjectID(),
		Username:     body.Username,
		Email:        body.Email,
		PasswordHash: hash,
		CreatedAt:    time.Now().UTC().String(),
		AppRole:      string(role),
		Reports:      []primitive.ObjectID{},
	}

	userID, err := u.Service.CreateUser(newUser)
	if err != nil {
		u.Logger.LogError(
			fmt.Sprintf("Error happened: %v", err),
		)
		c.IndentedJSON(500, gin.H{
			"message": "Internal server error",
		})
	}

	u.Logger.LogInfo(
		fmt.Sprintf("New user %s added", userID),
	)
	c.IndentedJSON(200, gin.H{
		"message": fmt.Sprintf("New user %s was succesfully created", userID),
	})
}

// POST /login
func (u *userRouter) LoginUser(c *gin.Context) {
	var body models.User

	if err := c.BindJSON(&body); err != nil {
		u.Logger.LogError(
			fmt.Sprintf("Error happened while binding JSON: %v", err),
		)
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "Error in handling request",
		})
		return
	}

	existingUser, err := u.Service.CheckExistingUser(body.Username)
	if err != nil {
		u.Logger.LogError("No user found")
		c.IndentedJSON(401, gin.H{
			"message": "No user found",
		})
		return
	}

	isCorrectPassword := u.Service.CheckPassword(existingUser.PasswordHash, body.PasswordHash)

	if !isCorrectPassword {
		u.Logger.LogInfo("Incorrect password")
		c.IndentedJSON(401, gin.H{
			"message": "Incorrect password",
		})
		return
	}

	token, err := u.Service.CreateToken(existingUser)

	if err != nil {
		u.Logger.LogError(
			fmt.Sprintf("Error happened while creating token: %v", err),
		)
		c.IndentedJSON(500, gin.H{
			"message": "Internal server error",
		})
		return
	}

	c.IndentedJSON(200, gin.H{
		"message": "Login succesful",
		"token":   token.Token,
	})
}
