package services

import (
	"errors"
	"main/pkg/middleware"
	"main/pkg/models"
	"main/pkg/repository"
)

type UserService interface {
	CreateUser(user models.User) (string, error)
	CheckExistingUser(username string) (models.User, error)
	DetermineRole(role string) (models.Role, error)
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

func (u *userService) CreateUser(user models.User) (string, error) {
	return u.Repository.Create(&user)
}

func (u *userService) CheckExistingUser(username string) (models.User, error) {
	r := u.Repository
	user, err := r.GetSingleUser(username)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (u *userService) DetermineRole(role string) (models.Role, error) {
	switch role {
	case "admin":
		return models.Admin, nil
	case "maintainer":
		return models.Maintainer, nil
	case "creator":
		return models.Creator, nil
	case "guest":
		return models.Guest, nil
	default:
		return models.Undefined, errors.New("Role undefined: " + role)
	}
}
