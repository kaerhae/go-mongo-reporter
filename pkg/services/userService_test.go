package services_test

import (
	"main/pkg/helpers"
	"main/pkg/middleware"
	"main/pkg/models"
	"main/pkg/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func initUserService() services.UserService {
	repo := helpers.InitMockUserRepository()
	return services.NewUserService(repo, middleware.NewSyslogger(false))
}

func TestRegistrationShouldReturnId(t *testing.T) {
	newUser := models.User{
		Username:     "test",
		Email:        "test@test.com",
		PasswordHash: "test",
		AppRole:      "guest",
		CreatedAt:    "nil",
	}
	s := initUserService()
	tUser, err := s.CreateUser(newUser)
	if err != nil {
		t.Fail()
	}

	assert.Equal(t, len(tUser), 24)
}
