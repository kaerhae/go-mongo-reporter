package helpers

import (
	"log"
	"main/pkg/models"
	"main/pkg/utils"
)

/*
*
*
*  FOR USER TESTS
*
*
 */

type mockUserRepository struct{}

type MockUserRepository interface {
	Create(user *models.User) (string, error)
	Get() ([]models.User, error)
	GetSingleUser(username string) (models.User, error)
	UpdateSingleUser(user models.User) error
	DeleteSingleUser(ID string) (int64, error)
}

func InitMockUserRepository() MockUserRepository {
	return &mockUserRepository{}
}

func (m *mockUserRepository) Create(user *models.User) (string, error) {
	return user.ID.Hex(), nil
}

// Get implements MockUserRepository.
func (m *mockUserRepository) Get() ([]models.User, error) {
	return []models.User{}, nil
}

func (m *mockUserRepository) GetSingleUser(username string) (models.User, error) {
	pwd, err := utils.HashPwd("strong-password")
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

// UpdateUser implements MockUserRepository.
func (m *mockUserRepository) UpdateSingleUser(user models.User) error {
	return nil
}
func (m *mockUserRepository) DeleteSingleUser(_ string) (int64, error) { return 1, nil }
