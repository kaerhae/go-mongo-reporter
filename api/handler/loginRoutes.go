package handler

import (
	"fmt"
	"main/pkg/models"
	"main/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
GET /login route. Requires admin rights. Takes LoginUser model as request body and validates body.
Then checks if user exists, if it exists checks if password matches to hashed password on db.

If no user, return 401 error. If invalid password return 401 error.

Finally creates a token and returns response with success message, userID, and token.
*/
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
