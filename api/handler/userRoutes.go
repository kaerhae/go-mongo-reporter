package handler

import (
	"fmt"
	"main/pkg/middleware"
	"main/pkg/models"
	"main/pkg/services"
	"main/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserRouter interface {
	Get(c *gin.Context)
	GetByID(c *gin.Context)
	LoginUser(c *gin.Context)
	PostNewUser(c *gin.Context)
	PostNewGuestUser(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
}

type userRouter struct {
	Service services.UserService
	Logger  middleware.Logger
}

func NewUserHandler(service services.UserService, logger middleware.Logger) UserRouter {
	return &userRouter{
		Service: service,
		Logger:  logger,
	}
}

func (u *userRouter) Get(c *gin.Context) {
	users, err := u.Service.GetAll()
	if err != nil {
		u.Logger.LogError(
			fmt.Sprintf("Error happened while fetching reports: %v", err),
		)
		c.IndentedJSON(500, gin.H{"message": fmt.Sprintf("Internal server error: %v", err)})
		return
	}

	c.IndentedJSON(200, users)
}

// GetById implements ReportRouter.
func (u *userRouter) GetByID(c *gin.Context) {
	id := c.Param("id")
	report, err := u.Service.GetByID(id)
	if err != nil {
		u.Logger.LogError(
			fmt.Sprintf("Error happened while fetching single user: %v", err),
		)
		c.IndentedJSON(400, gin.H{
			"message": fmt.Sprintf("Error: %v", err),
		})
		return
	}

	c.IndentedJSON(200, report)
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
		"token":   token,
	})
}

// POST /users
//
//nolint:dupl
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

	userID, err := u.Service.CreateUser(body)
	if err != nil {
		u.Logger.LogError(
			fmt.Sprintf("Error while creating user: %v", err),
		)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
		})
	}
	c.IndentedJSON(201, gin.H{
		"message": fmt.Sprintf("New user %s was succesfully created", userID),
	})
}

func (u *userRouter) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var body models.User
	err := c.BindJSON(&body)
	if err != nil {
		u.Logger.LogError(
			fmt.Sprintf("Error happened while binding JSON: %v", err),
		)
		c.IndentedJSON(400, gin.H{"message": "error on parsing request body"})
		return
	}

	existingUser, err := u.Service.GetByID(id)
	if err != nil {
		u.Logger.LogInfo("No report found")
		c.IndentedJSON(400, gin.H{"message": "No report found"})
		return
	}

	newReport := models.User{
		ID:           existingUser.ID,
		Username:     body.Username,
		Email:        body.Email,
		PasswordHash: existingUser.PasswordHash,
		CreatedAt:    existingUser.CreatedAt,
		Permission:   body.Permission,
	}
	err = u.Service.UpdateUser(newReport)
	if err != nil {
		u.Logger.LogError(
			fmt.Sprintf("Error happened while updating user: %v", err),
		)
		c.IndentedJSON(400, gin.H{"message": fmt.Sprintf("Internal server error: %v", err)})
		return
	}

	c.IndentedJSON(200, gin.H{
		"message": fmt.Sprintf("User \"%s\" was succesfully updated", body.Username),
	})
}

func (u *userRouter) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	deletedCount, err := u.Service.DeleteUser(id)
	if err != nil {
		u.Logger.LogError(
			fmt.Sprintf("Error happened while deleting user: %v", err),
		)
		c.IndentedJSON(500, gin.H{"message": "Internal server error"})
		return
	}
	u.Logger.LogInfo(fmt.Sprintf("Deleted %d users", deletedCount))
	c.IndentedJSON(200, gin.H{
		"message": fmt.Sprintf("Deleted %d users", deletedCount),
	})
}
