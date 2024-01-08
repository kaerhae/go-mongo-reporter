package routes

import (
	"fmt"
	"main/cmd/models"
	"main/cmd/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var validate = validator.New()

// POST /users
func PostNewUser(c *gin.Context) {
	var body models.User

	/* Bindataan request body muuttujaan body */
	if err := c.BindJSON(&body); err != nil {
		c.IndentedJSON(http.StatusBadRequest, "Error in handling request")
		return
	}

	_, err := services.CheckExistingUser(body.Username)

	/* Tarkistetaan löytyykö käyttäjää ennestään */
	if err == nil {
		c.IndentedJSON(http.StatusBadRequest, "Username already exists")
		return
	}

	if validationErr := validate.Struct(&body); validationErr != nil {
		c.IndentedJSON(http.StatusBadRequest, "Malformatted body")
		return
	}

	role, err := services.DetermineRole(string(body.App_Role))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, "Malformatted role")
		return
	}

	hash, err := services.HashPwd(body.Password_hash)
	if err != nil {
		c.IndentedJSON(500, "Server failed")
		return
	}

	newUser := models.User{
		ID:            primitive.NewObjectID(),
		Username:      body.Username,
		Email:         body.Email,
		Password_hash: hash,
		Created_At:    time.Now().UTC().String(),
		App_Role:      string(role),
	}
	u, err := services.CreateUser(newUser)
	if err != nil {
		c.IndentedJSON(500, gin.H{
			"message": "Internal server error",
		})
	}

	c.IndentedJSON(200, gin.H{
		"message": fmt.Sprintf("New user %s was succesfully created", u.Username),
	})

}

// POST /login
func LoginUser(c *gin.Context) {
	var body models.User

	if err := c.BindJSON(&body); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"message": "Error in handling request",
		})
		return
	}

	existingUser, err := services.CheckExistingUser(body.Username)
	if err != nil {
		c.IndentedJSON(401, gin.H{
			"message": "No user found",
		})
		return
	}

	isCorrectPassword := services.CheckPassword(existingUser.Password_hash, body.Password_hash)

	if !isCorrectPassword {
		c.IndentedJSON(401, gin.H{
			"message": "Incorrect password",
		})
		return
	}

	token, err := services.CreateToken(body.Username)

	if err != nil {
		c.IndentedJSON(500, gin.H{
			"message": "Internal server error",
		})
		return
	}

	c.IndentedJSON(200, gin.H{
		"message": "Login succesfull",
		"token":   token.Token,
	})
}
