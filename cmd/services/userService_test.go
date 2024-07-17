package services_test

import (
	"main/cmd/helpers"
	"main/cmd/models"
	"main/cmd/services"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegistrationShouldReturnId(t *testing.T) {
	newUser := models.User{
		Username:     "test",
		Email:        "test@test.com",
		PasswordHash: "test",
		AppRole:      "guest",
		CreatedAt:    "nil",
	}
	repo := helpers.InitMockRepository()
	s := services.NewUserService(repo)
	tUser, err := s.CreateUser(newUser)
	if err != nil {
		t.Fail()
	}

	assert.Equal(t, len(tUser), 24)
}

func TestRoleDeterminationShouldReturnCorrect(t *testing.T) {
	repo := helpers.InitMockRepository()
	s := services.NewUserService(repo)
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
	repo := helpers.InitMockRepository()
	s := services.NewUserService(repo)
	g, err := s.DetermineRole("other")
	if err == nil || g != models.Undefined {
		t.Fail()
	}
}
