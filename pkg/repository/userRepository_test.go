package repository_test

import (
	"main/pkg/helpers"
	"main/pkg/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestRepoCreateShouldReturnID(t *testing.T) {
	repo := helpers.InitMockRepository()
	id := primitive.NewObjectID()
	newUser := models.User{
		ID:           id,
		Username:     "test",
		Email:        "test@test.com",
		PasswordHash: "pass",
		AppRole:      "guest",
		CreatedAt:    "nil",
	}
	u, err := repo.Create(&newUser)

	if err != nil {
		t.Fail()
	}

	assert.Equal(t, id.Hex(), u)
}

func TestRepoSingleUserShouldBeRetrieved(t *testing.T) {
	repo := helpers.InitMockRepository()

	user, err := repo.GetSingleUser("user")
	if err != nil {
		t.Fail()
	}
	assert.Equal(t, "user", user.Username)
	assert.Equal(t, "test@test.com", user.Email)
	assert.Equal(t, "guest", user.AppRole)
	assert.Equal(t, "2023-01-01", user.CreatedAt)
}
