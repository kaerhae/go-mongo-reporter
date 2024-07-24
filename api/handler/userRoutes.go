package handler

import (
	"fmt"
	"main/pkg/middleware"
	"main/pkg/models"
	"main/pkg/services"
	"main/pkg/utils"
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
	var body models.CreateUser
	if err := c.BindJSON(&body); err != nil {
		u.Logger.LogError(
			fmt.Sprintf("Error happened while binding JSON: %v", err),
		)
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "Error in handling request",
		})
		return
	}

	if body.Username == "" || body.Password == "" || body.Email == "" {
		u.Logger.LogError(fmt.Sprintf("Malformatted body: %v\n", body))
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "Malformatted body",
		})
		return
	}

	_, err := u.Service.CheckExistingUser(body.Username)

	// check if user exists already
	if err == nil {

		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "Username already exists",
		})
		return
	}

	if validationErr := validate.Struct(&body); validationErr != nil {
		u.Logger.LogError(
			fmt.Sprintf("Error happened while validating body: %v", err),
		)
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "Malformatted body",
		})
		return
	}
	role, err := u.Service.DetermineRole(string(body.AppRole))
	if err != nil {
		u.Logger.LogError(
			fmt.Sprintf("Malformatted appRole: %v", err),
		)
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "Malformatted role",
		})
		return
	}
	hash, err := utils.HashPwd(body.Password)
	if err != nil {
		u.Logger.LogError(
			fmt.Sprintf("Error happened while hashing password: %v", err),
		)
		c.IndentedJSON(500, gin.H{
			"message": "Server failed",
		})
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
		return
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
	var body models.LoginUser

	if err := c.BindJSON(&body); err != nil {
		u.Logger.LogError(
			fmt.Sprintf("Error happened while binding JSON: %v", err),
		)
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "Error in handling request",
		})
		return
	}

	if body.Username == "" || body.Password == "" {
		u.Logger.LogError("Malformatted body")
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "Malformatted body",
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

	err = utils.CheckPassword(existingUser.PasswordHash, body.Password)

	if err != nil {
		u.Logger.LogInfo(fmt.Sprintf("Error while checking password: %v", err))
		c.IndentedJSON(401, gin.H{
			"message": "Incorrect password",
		})
		return
	}

	token, err := utils.CreateToken(existingUser)

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
		"userID":  existingUser.ID.Hex(),
		"token":   token.Token,
	})
}
