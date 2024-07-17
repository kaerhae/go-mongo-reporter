package helpers

import (
	"log"
	"main/cmd/models"

	"golang.org/x/crypto/bcrypt"
)

type mockUserRepository struct{}
type MockUserRepository interface {
	Create(user *models.User) (string, error)
	GetSingleUser(username string) (models.User, error)
}

func InitMockRepository() MockUserRepository {
	return &mockUserRepository{}
}

func (m *mockUserRepository) Create(user *models.User) (string, error) {
	return user.ID.Hex(), nil
}

func (m *mockUserRepository) GetSingleUser(username string) (models.User, error) {
	pwd, err := TestingEnvHashPwd("strong-password")
	if err != nil {
		log.Fatal(err)
	}
	return models.User{
		Username:     username,
		PasswordHash: pwd,
		Email:        "test@test.com",
		AppRole:      "guest",
		CreatedAt:    "2023-01-01",
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
