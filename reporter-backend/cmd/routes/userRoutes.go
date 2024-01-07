package routes

import (
	"log"
	"main/reporter-backend/cmd/models"
	"main/reporter-backend/cmd/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
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

	isExisting := services.CheckExistingUser(body.Username)

	/* Tarkistetaan löytyykö käyttäjää ennestään */
	if isExisting {
		c.IndentedJSON(http.StatusBadRequest, "Username already exists")
		return
	}

	if validationErr := validate.Struct(&body); validationErr != nil {
		c.IndentedJSON(http.StatusBadRequest, "Malformatted body")
		return
	}

	role, err := models.DetermineRole(string(body.App_Role))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, "Malformatted role")
		return
	}

	hash, err := hashPwd(body.Password_hash)
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
	services.CreateUser(newUser)
}

// POST /login
func loginUser(c *gin.Context) {
	var body models.User

	if err := c.BindJSON(&body); err != nil {
		c.IndentedJSON(http.StatusBadRequest, "Error in handling request")
		return
	}

}

func hashPwd(password string) (string, error) {
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

func checkPassword(hashedPassword string, plainPassword string) bool {
	bytePlain := []byte(plainPassword)
	byteHashed := []byte(hashedPassword)

	err := bcrypt.CompareHashAndPassword(byteHashed, bytePlain)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
