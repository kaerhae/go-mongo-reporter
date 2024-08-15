package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"main/pkg/helpers"
	"main/pkg/middleware"
	"main/pkg/models"
	"main/pkg/services"
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

func (s *MockUserService) GetAll() ([]models.User, error) {
	return []models.User{
		{Username: "test", Email: "test@test.com", Permission: models.Permission{}},
		{Username: "test2", Email: "test@test.com", Permission: models.Permission{}},
	}, nil
}

func (s *MockUserService) GetByID(id string) (models.User, error) {
	return models.User{Username: "test", Email: "test@test.com", Permission: models.Permission{}}, nil
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
func (s *MockUserService) CreateUser(user models.CreateUser) (string, error) {
	return "1234", nil
}

// CreateGuestUser implements services.UserService.
func (s *MockUserService) CreateGuestUser(user models.CreateGuestUser) (string, error) {
	return "1234", nil
}

func (s *MockUserService) UpdateUser(user models.User) error {
	return nil
}

func (s *MockUserService) DeleteUser(id string) (int64, error) {
	return 0, nil
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

func TestLoginUserShouldBeSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	os.Setenv("DATABASE_URI", "test")
	/* This test is going to compare user object password to object that MockUserRepository GetSingleUser retrieves */
	user := models.LoginUser{
		Username: "testerUser",
		Password: "strong-password",
	}

	repo := helpers.InitMockUserRepository()
	s := services.NewUserService(repo, middleware.NewSyslogger(false))
	handler := NewUserHandler(s, middleware.NewSyslogger(false))
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
	gin.SetMode(gin.TestMode)
	/* This test is going to compare user object password to object that MockUserRepository GetSingleUser retrieves */
	user := models.LoginUser{
		Username: "testerUser",
		Password: "weak-wrong-password",
	}

	repo := helpers.InitMockUserRepository()
	s := services.NewUserService(repo, middleware.NewSyslogger(false))
	handler := NewUserHandler(s, middleware.NewSyslogger(false))
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

func TestGetAll(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	repo := helpers.InitMockUserRepository()
	s := &MockUserService{Repository: repo}
	l := middleware.NewSyslogger(false)
	userHandler := NewUserHandler(s, l)
	r.GET("/user-management/users", userHandler.Get)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/user-management/users", nil)
	r.ServeHTTP(w, req)

	var users []models.User
	json.Unmarshal(w.Body.Bytes(), &users)

	assert.Equal(t, 200, w.Code)
	assert.NotEmpty(t, users)
}

func TestPostNewUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	newUser := models.CreateUser{
		Username: "testerUser",
		Email:    "test@test.com",
		Password: "passhash",
		Permission: models.Permission{
			Admin: true,
			Write: true,
			Read:  true,
		},
	}

	repo := helpers.InitMockUserRepository()
	s := &MockUserService{Repository: repo}
	handler := NewUserHandler(s, middleware.NewSyslogger(false))
	router := SetUpRouter()
	router.POST("/user-management/users", handler.PostNewUser)
	w := httptest.NewRecorder()
	payload, _ := json.Marshal(newUser)

	req, _ := http.NewRequest("POST", "/user-management/users", bytes.NewBuffer(payload))

	router.ServeHTTP(w, req)

	var response SingleMessageResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Equal(t, 201, w.Code)
	assert.Equal(t, "New user 1234 was succesfully created", response.Message)
}

func TestUpdateUser_ShouldUpdate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	user := models.User{
		Username:   "test2",
		Email:      "test@test2.com",
		Permission: models.Permission{},
	}
	repo := helpers.InitMockUserRepository()
	s := &MockUserService{Repository: repo}
	l := middleware.NewSyslogger(false)
	userHandler := NewUserHandler(s, l)
	r.PUT("/user-management/users/:id", userHandler.UpdateUser)
	w := httptest.NewRecorder()
	payload, _ := json.Marshal(user)
	req, _ := http.NewRequest(http.MethodPut, "/user-management/users/1", bytes.NewBuffer(payload))
	r.ServeHTTP(w, req)

	var users []models.User
	json.Unmarshal(w.Body.Bytes(), &users)

	assert.Equal(t, 200, w.Code)
}

func TestDeleteUser_ShouldDelet(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	repo := helpers.InitMockUserRepository()
	s := &MockUserService{Repository: repo}
	l := middleware.NewSyslogger(false)
	userHandler := NewUserHandler(s, l)
	r.DELETE("/user-management/users/:id", userHandler.DeleteUser)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodDelete, "/user-management/users/1", nil)
	r.ServeHTTP(w, req)

	var users []models.User
	json.Unmarshal(w.Body.Bytes(), &users)

	assert.Equal(t, 200, w.Code)
}
