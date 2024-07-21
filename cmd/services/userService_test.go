package services_test

import (
	"main/cmd/helpers"
	"main/cmd/middleware"
	"main/cmd/models"
	"main/cmd/services"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func initUserService() services.UserService {
	repo := helpers.InitMockRepository()
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

func TestRoleDeterminationShouldReturnCorrect(t *testing.T) {
	s := initUserService()
	g, err := s.DetermineRole("guest")
	if err != nil {
		t.Fail()
	}
	a, err := s.DetermineRole("admin")
	if err != nil {
		t.Fail()
	}
	c, err := s.DetermineRole("creator")
	if err != nil || g != models.Guest || a != models.Admin || c != models.Creator {
		t.Fail()
	}
}

func TestRoleDeterminationShouldRetunUndefined(t *testing.T) {
	s := initUserService()
	g, err := s.DetermineRole("other")
	if err == nil || g != models.Undefined {
		t.Fail()
	}
}

func TestTokenCreation(t *testing.T) {
	os.Setenv("SECRET_KEY", "123")
	user := models.User{
		Username: "tester",
	}

	s := initUserService()

	claims, err := s.CreateToken(user)
	if err != nil {
		t.Fail()
	}
	assert.Equal(t, user.Username, claims.Username)
	assert.Equal(t, user.AppRole, claims.AppRole)
	assert.NotEmpty(t, claims.Token)

	os.Unsetenv("SECRET_KEY")
}
