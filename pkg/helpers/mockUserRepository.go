package helpers

import (
	"log"
	"main/pkg/models"
	"main/pkg/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
*
*
*  FOR USER TESTS
*
*
 */

type mockUserRepository struct{}

// GetSingleUserById implements MockUserRepository.
func (m *mockUserRepository) GetSingleUserByID(id string) (models.User, error) {
	i, _ := primitive.ObjectIDFromHex(id)

	return models.User{
		ID:           i,
		Username:     "test",
		Email:        "test@test.com",
		PasswordHash: "1234",
		Permission:   models.Permission{}}, nil
}

type MockUserRepository interface {
	Create(user *models.User) (string, error)
	Get() ([]models.User, error)
	GetSingleUserByID(id string) (models.User, error)
	GetSingleUserByUsername(id string) (models.User, error)
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
	p := models.Permission{}
	p.SetDefaultPermissions()
	list := []models.User{
		{Username: "test", Email: "test@test.com", Permission: p},
		{Username: "test2", Email: "test@test.com", Permission: p},
	}
	return list, nil
}

func (m *mockUserRepository) GetSingleUserByUsername(id string) (models.User, error) {
	pwd, err := utils.HashPwd("strong-password")
	if err != nil {
		log.Fatal(err)
	}
	p := models.Permission{}
	p.SetDefaultPermissions()
	i, _ := primitive.ObjectIDFromHex(id)
	return models.User{
		ID:           i,
		Username:     "test",
		PasswordHash: pwd,
		Email:        "test@test.com",
		Permission:   p,
		CreatedAt:    "2023-01-01",
	}, nil
}

// UpdateUser implements MockUserRepository.
func (m *mockUserRepository) UpdateSingleUser(_ models.User) error {
	return nil
}
func (m *mockUserRepository) DeleteSingleUser(_ string) (int64, error) { return 1, nil }
