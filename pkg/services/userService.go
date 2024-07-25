package services

import (
	"main/pkg/middleware"
	"main/pkg/models"
	"main/pkg/repository"
)

type UserService interface {
	CreateUser(user models.User) (string, error)
	CheckExistingUser(username string) (models.User, error)
	GetAll() ([]models.User, error)
	UpdateUser(user models.User) error
	DeleteUser(id string) (int64, error)
}

type userService struct {
	Repository repository.UserRepository
	Logger     middleware.Logger
}

func NewUserService(repo repository.UserRepository, logger middleware.Logger) UserService {
	return &userService{
		Repository: repo,
		Logger:     logger,
	}
}

func (u *userService) GetAll() ([]models.User, error) {
	return u.Repository.Get()
}

func (u *userService) CreateUser(user models.User) (string, error) {
	return u.Repository.Create(&user)
}

func (u *userService) CheckExistingUser(username string) (models.User, error) {
	return u.Repository.GetSingleUser(username)
}

func (u *userService) UpdateUser(user models.User) error {
	return u.Repository.UpdateSingleUser(user)
}

func (u *userService) DeleteUser(id string) (int64, error) {
	return u.Repository.DeleteSingleUser(id)
}
