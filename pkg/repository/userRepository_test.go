package repository_test

import (
	"main/pkg/helpers"
	"main/pkg/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestRepoCreateShouldReturnID(t *testing.T) {
	repo := helpers.InitMockUserRepository()
	id := primitive.NewObjectID()
	p := models.Permission{}
	p.SetDefaultPermissions()
	newUser := models.User{
		ID:           id,
		Username:     "test",
		Email:        "test@test.com",
		PasswordHash: "pass",
		CreatedAt:    "nil",
		Permission:   p,
	}
	u, err := repo.Create(&newUser)

	if err != nil {
		t.Fail()
	}

	assert.Equal(t, id.Hex(), u)
}

func TestRepoSingleUserShouldBeRetrieved(t *testing.T) {
	repo := helpers.InitMockUserRepository()
	p := models.Permission{}
	p.SetDefaultPermissions()
	user, err := repo.GetSingleUserByUsername("test")
	if err != nil {
		t.Fail()
	}
	assert.Equal(t, "test", user.Username)
	assert.Equal(t, "test@test.com", user.Email)
	assert.Equal(t, p, user.Permission)
	assert.Equal(t, "2023-01-01", user.CreatedAt)
}

func TestRepoCreate_ShouldCreateUser(t *testing.T) {
	repo := helpers.InitMockUserRepository()
	id, err := primitive.ObjectIDFromHex("123456789012345678901234")
	if err != nil {
		t.Fail()
	}
	p := models.Permission{}
	p.SetDefaultPermissions()
	newUser := models.User{
		ID:         id,
		Username:   "test",
		Email:      "test@test.com",
		Permission: p,
	}
	objectId, err := repo.Create(&newUser)
	if err != nil {
		t.Fail()
	}

	assert.Equal(t, objectId, id.Hex())
}
