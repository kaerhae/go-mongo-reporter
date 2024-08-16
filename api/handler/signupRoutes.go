package handler

import (
	"fmt"
	"main/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
POST /signup route. Allowed all access. Takes CreateGuestUser model as request body and validates body.

Checks that user does not exist. If exists, return 409 error.

Finally calls CreateGuestUser method and if successful, returns response with success message.
*/
//
//nolint:dupl
func (u *userRouter) PostNewGuestUser(c *gin.Context) {
	/* CreateUser model without permission field */
	var body models.CreateGuestUser
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
		u.Logger.LogError(fmt.Sprintf("Malformatted guest user body: %v\n", body))
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "Malformatted guest user body",
		})
		return
	}

	_, err := u.Service.CheckExistingUser(body.Username)

	// check if user exists already
	if err == nil {
		c.IndentedJSON(http.StatusConflict, gin.H{
			"message": "Username already exists",
		})
		return
	}

	userID, err := u.Service.CreateGuestUser(body)
	if err != nil {
		u.Logger.LogError(
			fmt.Sprintf("Error while creating guest user: %v", err),
		)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
		})
	}
	c.IndentedJSON(200, gin.H{
		"message": fmt.Sprintf("New user %s was succesfully created", userID),
	})
}
