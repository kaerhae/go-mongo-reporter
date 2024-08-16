package handler

import (
	"fmt"
	"main/pkg/middleware"
	"main/pkg/models"
	"main/pkg/services"
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
	UpdatePassword(c *gin.Context)
	UpdateUserPermissions(c *gin.Context)
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

/*
GET /user-management/users route. Requires admin rights. Returns all users from database.
*/
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

/*
GET /user-management/users/:id route. Requires admin rights. Returns single user by ID from database.
*/
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

/*
POST /user-management/users route. Requires admin rights. Takes CreateUser model as request body and validates body.
Then checks if user exists, if exists, returns 400 error.

Finally creates a new user and returns response with success message.
*/
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

/*
PUT /user-management/users route. Requires admin rights. Takes User model as request body and validates body.
Then checks if user exists, if not, returns 400 error.

Finally creates a updated user object to db and returns response with success message.
*/
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

	/* UpdateUser has allowed update properties for username and email */
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

/*
PUT /change-password route. Admin and selective non-admin rights. Takes UserPassword model as request body and validates body.
Also fetches userID and admin info from gin context. If request is not from admin, checks that user is same. If not return 403 error.

Finally calls UpdatePassword method and if successful, returns response with success message.
*/
func (u *userRouter) UpdatePassword(c *gin.Context) {
	userID, isAdmin, err := middleware.GetSessionData(c)
	if err != nil {
		u.Logger.LogError(
			fmt.Sprintf("Error happened while getting session data: %v", err),
		)
		c.IndentedJSON(500, gin.H{"message": "internal server error"})
		return
	}
	var body models.UserPasswordChange
	err = c.BindJSON(&body)
	if err != nil {
		u.Logger.LogError(
			fmt.Sprintf("Error happened while binding JSON: %v", err),
		)
		c.IndentedJSON(400, gin.H{"message": "error on parsing request body"})
		return
	}

	if body.UserID == "" || body.NewPassword == "" || body.OldPassword == "" {
		u.Logger.LogError("Invalid body")
		c.IndentedJSON(400, gin.H{"message": "Invalid body"})
		return
	}

	if !isAdmin && userID != body.UserID {
		u.Logger.LogError("Non-admin user is trying to change other users password")
		c.IndentedJSON(http.StatusForbidden, gin.H{"message": "illegal request"})
	}

	// Call UpdatePassword method with userID of user, which password is going to be changed
	err = u.Service.UpdatePassword(body.UserID, body.OldPassword, body.NewPassword)
	if err != nil {
		u.Logger.LogError(
			fmt.Sprintf("Error happened while updating user password: %v", err),
		)
		c.IndentedJSON(500, gin.H{"message": fmt.Sprintf("error while updating password: %v", err)})
		return
	}

	c.IndentedJSON(200, gin.H{"message": "Password updated successfully"})
}

/*
PUT /user-management/change-permissions route. Requires admin rights. Takes UserPermissionUpdate model as request body and validates body.

Finally calls UpdateUserPermissions method and if successful, returns response with success message.
*/
func (u *userRouter) UpdateUserPermissions(c *gin.Context) {
	var body models.UserPermissionUpdate
	err := c.BindJSON(&body)
	if err != nil {
		u.Logger.LogError(
			fmt.Sprintf("Error happened while binding JSON: %v", err),
		)
		c.IndentedJSON(400, gin.H{"message": "error on parsing request body"})
		return
	}

	if body.UserID == "" {
		u.Logger.LogError("Invalid body")
		c.IndentedJSON(400, gin.H{"message": "Invalid body"})
		return
	}

	err = u.Service.UpdateUserPermissions(body.UserID, body.Permissions)
	if err != nil {
		u.Logger.LogError(
			fmt.Sprintf("Error happened while updating permissions: %v", err),
		)
		c.IndentedJSON(500, gin.H{"message": "internal server error"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "permissions updated"})

}

/*
DELETE /user-management/users route. Requires admin rights. Takes id from url parameters.

Calls DeleteUser method and if successful, returns response with success message.
*/
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
