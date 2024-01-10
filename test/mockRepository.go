package test

import (
	"log"
	"main/cmd/models"

	"golang.org/x/crypto/bcrypt"
)

type mockUserRepository struct{}

func (m *mockUserRepository) Create(user *models.User) (string, error) {
	return user.ID.Hex(), nil
}

func (m *mockUserRepository) GetSingleUser(username string) (models.User, error) {
	pwd, err := TestingEnvHashPwd("strong-password")
	if err != nil {
		log.Fatal(err)
	}
	return models.User{
		Username:      "test",
		Password_hash: pwd,
		Email:         "test@test.com",
		App_Role:      "guest",
		Created_At:    "2023-01-01",
	}, nil
}

func TestingEnvHashPwd(pass string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(pass),
		bcrypt.MinCost,
	)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return string(hash), nil
}
