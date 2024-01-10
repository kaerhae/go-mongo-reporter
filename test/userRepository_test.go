package test

import (
	"main/cmd/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepoCreateShouldReturnID(t *testing.T) {
	repo := &mockUserRepository{}
	newUser := models.User{
		ID:            id,
		Username:      "test",
		Email:         "test@test.com",
		Password_hash: "pass",
		App_Role:      "guest",
		Created_At:    "nil",
	}
	u, err := repo.Create(&newUser)

	if err != nil {
		t.Fail()
	}

	assert.Equal(t, id.Hex(), u)
}

func TestRepoSingleUserShouldBeRetrieved(t *testing.T) {
	repo := &mockUserRepository{}

	user, err := repo.GetSingleUser("user")
	if err != nil {
		t.Fail()
	}
	assert.Equal(t, "test", user.Username)
	assert.Equal(t, "test@test.com", user.Email)
	assert.Equal(t, "guest", user.App_Role)
	assert.Equal(t, "2023-01-01", user.Created_At)
}
