package handler

import (
	"bytes"
	"encoding/json"
	"main/pkg/helpers"
	"main/pkg/middleware"
	"main/pkg/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPostNewGuestUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	newUser := models.CreateUser{
		Username: "testerUser",
		Email:    "test@test.com",
		Password: "passhash",
	}

	repo := helpers.InitMockUserRepository()
	s := &MockUserService{Repository: repo}
	handler := NewUserHandler(s, middleware.NewSyslogger(false))
	router := SetUpRouter()
	router.POST("/signup", handler.PostNewGuestUser)
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
