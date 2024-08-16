package helpers

import (
	"errors"
	"main/pkg/middleware"
	"main/pkg/models"
)

type MockUserService struct {
	Repository MockUserRepository
	Logger     middleware.Logger
}

// UpdateUserPermissions implements services.UserService.
func (s *MockUserService) UpdateUserPermissions(_ string, _ models.Permission) error {
	panic("unimplemented")
}

func (s *MockUserService) GetAll() ([]models.User, error) {
	return []models.User{
		{Username: "test", Email: "test@test.com", Permission: models.Permission{}},
		{Username: "test2", Email: "test@test.com", Permission: models.Permission{}},
	}, nil
}

func (s *MockUserService) GetByID(_ string) (models.User, error) {
	return models.User{Username: "test", Email: "test@test.com", Permission: models.Permission{}}, nil
}

// CheckExistingUser implements services.UserService.
func (s *MockUserService) CheckExistingUser(_ string) (models.User, error) {
	return models.User{}, errors.New("")
}

// CheckPassword implements services.UserService.
func (s *MockUserService) CheckPassword(hashedPassword string, plainPassword string) bool {
	return true
}

// CreateToken implements services.UserService.
func (s *MockUserService) CreateToken(_ models.User) (*models.Claims, error) {
	return &models.Claims{}, nil
}

// CreateUser implements services.UserService.
func (s *MockUserService) CreateUser(_ models.CreateUser) (string, error) {
	return "1234", nil
}

// CreateGuestUser implements services.UserService.
func (s *MockUserService) CreateGuestUser(_ models.CreateGuestUser) (string, error) {
	return "1234", nil
}

func (s *MockUserService) UpdateUser(_ models.User) error {
	return nil
}

// UpdatePassword implements services.UserService.
func (s *MockUserService) UpdatePassword(_ string, _ string, _ string) error {
	return nil
}

func (s *MockUserService) DeleteUser(id string) (int64, error) {
	return 0, nil
}

// HashPwd implements services.UserService.
func (s *MockUserService) HashPwd(_ string) (string, error) {
	return "", nil
}

type SingleMessageResponse struct {
	Message string
}

type LoginResponse struct {
	Message string
	Token   string
}
