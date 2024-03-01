package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"main/cmd/models"
	"main/cmd/routes"
	"main/cmd/services"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var id primitive.ObjectID = primitive.NewObjectID()

type MockService struct {
	Repository mockUserRepository
}

type SingleMessageResponse struct {
	Message string
}

type LoginResponse struct {
	Message string
	Token   string
}

// CheckExistingUser implements services.UserService.
func (s *MockService) CheckExistingUser(username string) (models.User, error) {
	return models.User{}, errors.New("")
}

// CheckPassword implements services.UserService.
func (s *MockService) CheckPassword(hashedPassword string, plainPassword string) bool {
	return true
}

// CreateToken implements services.UserService.
func (s *MockService) CreateToken(username string, app_role string) (*models.Claims, error) {
	return &models.Claims{}, nil
}

// CreateUser implements services.UserService.
func (s *MockService) CreateUser(user models.User) (string, error) {
	return "1234", nil
}

// DetermineRole implements services.UserService.
func (s *MockService) DetermineRole(role string) (models.Role, error) {
	return models.Guest, nil
}

// HashPwd implements services.UserService.
func (s *MockService) HashPwd(password string) (string, error) {
	return "", nil
}

func SetUpRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func TestPostNewUser(t *testing.T) {
	newUser := models.User{
		Username:      "testerUser",
		Email:         "test@test.com",
		Password_hash: "passhash",
		App_Role:      "guest",
		Created_At:    "nil",
	}

	repo := &mockUserRepository{}
	s := &MockService{Repository: *repo}
	handler := routes.NewUserHandler(s)
	router := SetUpRouter()
	router.POST("/signup", handler.PostNewUser)
	w := httptest.NewRecorder()
	payload, _ := json.Marshal(newUser)

	req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(payload))

	router.ServeHTTP(w, req)

	var response SingleMessageResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)

	assert.Nil(t, err)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "New user 1234 was succesfully created", response.Message)
}

func TestLoginUserShouldBeSuccess(t *testing.T) {
	/* This test is going to compare user object password to object that MockUserRepository GetSingleUser retrieves */
	user := models.User{
		Username:      "testerUser",
		Password_hash: "strong-password",
	}

	repo := &mockUserRepository{}
	s := services.NewUserService(repo)
	handler := routes.NewUserHandler(s)
	router := SetUpRouter()
	router.POST("/login", handler.LoginUser)
	w := httptest.NewRecorder()
	payload, _ := json.Marshal(user)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(payload))

	router.ServeHTTP(w, req)
	var response LoginResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)

	assert.Nil(t, err)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "Login succesful", response.Message)
}

func TestLoginUserShouldNotBeSuccess(t *testing.T) {
	/* This test is going to compare user object password to object that MockUserRepository GetSingleUser retrieves */
	user := models.User{
		Username:      "testerUser",
		Password_hash: "weak-wrong-password",
	}

	repo := &mockUserRepository{}
	s := services.NewUserService(repo)
	handler := routes.NewUserHandler(s)
	router := SetUpRouter()
	router.POST("/login", handler.LoginUser)
	w := httptest.NewRecorder()
	payload, _ := json.Marshal(user)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(payload))

	router.ServeHTTP(w, req)
	var response SingleMessageResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)

	assert.Nil(t, err)
	assert.Equal(t, 401, w.Code)
	assert.Equal(t, "Incorrect password", response.Message)
}
