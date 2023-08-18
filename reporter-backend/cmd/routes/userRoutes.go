package routes

import (
	"cmd/reporter-backend/cmd/models"
	"cmd/reporter-backend/cmd/service"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

func PostNewUser(c *gin.Context) {
	var body models.User

	/* Bindataan request body muuttujaan body */
	if err := c.BindJSON(&body); err != nil {
		c.IndentedJSON(http.StatusBadRequest, "Error in handling request")
	}

	isExisting := service.CheckExistingUser(body.Username)

	/* Tarkistetaan löytyykö käyttäjää ennestään */
	if isExisting == true {
		c.IndentedJSON(http.StatusBadRequest, "Username already exists")
	}

	if validationErr := validate.Struct(&body); validationErr != nil {
		c.IndentedJSON(http.StatusBadRequest, "Malformatted body")
	}

	role, err := models.DetermineRole(string(body.Role))
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, "Malformatted role")
	}

	newUser := models.User{
		ID:         primitive.NewObjectID(),
		Username:   body.Username,
		Email:      body.Email,
		Password:   hashPwd(body.Password),
		Created_At: time.Now().UTC().String(),
		Role:       string(role),
	}

	service.CreateUser(newUser)

}

func hashPwd(password string) string {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.MinCost,
	)
	if err != nil {
		log.Fatal(err)
	}

	return string(hash)
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
