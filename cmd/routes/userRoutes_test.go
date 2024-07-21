package routes

import (
	"bytes"
	"encoding/json"
	"errors"
	"main/cmd/helpers"
	"main/cmd/middleware"
	"main/cmd/models"
	"main/cmd/services"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type MockUserService struct {
	Repository helpers.MockUserRepository
	Logger     middleware.Logger
}

// CheckExistingUser implements services.UserService.
func (s *MockUserService) CheckExistingUser(username string) (models.User, error) {
	return models.User{}, errors.New("")
}

// CheckPassword implements services.UserService.
func (s *MockUserService) CheckPassword(hashedPassword string, plainPassword string) bool {
	return true
}

// CreateToken implements services.UserService.
func (s *MockUserService) CreateToken(user models.User) (*models.Claims, error) {
	return &models.Claims{}, nil
}

// CreateUser implements services.UserService.
func (s *MockUserService) CreateUser(user models.User) (string, error) {
	return "1234", nil
}

// DetermineRole implements services.UserService.
func (s *MockUserService) DetermineRole(role string) (models.Role, error) {
	return models.Guest, nil
}

// HashPwd implements services.UserService.
func (s *MockUserService) HashPwd(password string) (string, error) {
	return "", nil
}

type SingleMessageResponse struct {
	Message string
}

type LoginResponse struct {
	Message string
	Token   string
}

func SetUpRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	return router
}

func TestPostNewUser(t *testing.T) {
	newUser := models.User{
		Username:     "testerUser",
		Email:        "test@test.com",
		PasswordHash: "passhash",
		AppRole:      "guest",
		CreatedAt:    "nil",
	}

	repo := helpers.InitMockRepository()
	s := &MockUserService{Repository: repo}
	handler := NewUserHandler(s, middleware.NewSyslogger())
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
	os.Setenv("DATABASE_URI", "test")
	/* This test is going to compare user object password to object that MockUserRepository GetSingleUser retrieves */
	user := models.User{
		Username:     "testerUser",
		PasswordHash: "strong-password",
	}

	repo := helpers.InitMockRepository()
	s := services.NewUserService(repo, middleware.NewSyslogger())
	handler := NewUserHandler(s, middleware.NewSyslogger())
	router := gin.Default()
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
	os.Unsetenv("DATABASE_URI")
}

func TestLoginUserShouldNotBeSuccess(t *testing.T) {
	/* This test is going to compare user object password to object that MockUserRepository GetSingleUser retrieves */
	user := models.User{
		Username:     "testerUser",
		PasswordHash: "weak-wrong-password",
	}

	repo := helpers.InitMockRepository()
	s := services.NewUserService(repo, middleware.NewSyslogger())
	handler := NewUserHandler(s, middleware.NewSyslogger())
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
