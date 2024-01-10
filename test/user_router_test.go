package test

import (
	"context"
	"main/cmd/models"
	"main/cmd/services"
	"testing"

	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

type mockUserRepository struct{}

func (m *mockUserRepository) Create(ctx context.Context, user *models.User) (string, error) {
	return "", nil
}

func TestRoleDeterminationShouldReturnCorrect(t *testing.T) {
	g, err := services.DetermineRole("guest")
	if err != nil {
		t.Fail()
	}
	a, err := services.DetermineRole("admin")
	if err != nil {
		t.Fail()
	}
	c, err := services.DetermineRole("creator")
	if err != nil || g != models.Guest || a != models.Admin || c != models.Creator {
		t.Fail()
	}
}

func TestRoleDeterminationShouldRetunUndefined(t *testing.T) {
	g, err := services.DetermineRole("other")
	if err == nil || g != models.Undefined {
		t.Fail()
	}
}

func TestRegistration(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("test create", func(mt *mtest.T) {

		mt.AddMockResponses(mtest.CreateSuccessResponse())
		newUser := models.User{
			Username:      "test",
			Email:         "test@test.com",
			Password_hash: "test",
			App_Role:      "guest",
			Created_At:    "nil",
		}
		repo := &mockUserRepository{}
	})
}
