package services

import (
	"fmt"
	"main/pkg/middleware"
	"main/pkg/models"
	"main/pkg/repository"
	"main/pkg/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService interface {
	CreateUser(user models.CreateUser) (string, error)
	CreateGuestUser(user models.CreateGuestUser) (string, error)
	CheckExistingUser(username string) (models.User, error)
	GetByID(id string) (models.User, error)
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

/* User will be created with u.CreateUser, but this method makes sure that guest has correct permissions */
func (u *userService) CreateGuestUser(user models.CreateGuestUser) (string, error) {
	guestPermissions := models.Permission{
		Admin: false,
		Write: false,
		Read:  true,
	}
	insertableUser := models.CreateUser{
		Username:   user.Username,
		Email:      user.Email,
		Password:   user.Password,
		Permission: guestPermissions,
	}

	return u.CreateUser(insertableUser)
}

func (u *userService) CreateUser(user models.CreateUser) (string, error) {
	hash, err := utils.HashPwd(user.Password)
	if err != nil {
		u.Logger.LogError(
			fmt.Sprintf("Error happened while hashing password: %v", err),
		)
		return "", err
	}

	newUser := models.User{
		ID:           primitive.NewObjectID(),
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: hash,
		CreatedAt:    time.Now().UTC().String(),
		Permission:   user.Permission,
		Reports:      []primitive.ObjectID{},
	}

	return u.Repository.Create(&newUser)
}

func (u *userService) CheckExistingUser(username string) (models.User, error) {
	return u.Repository.GetSingleUserByUsername(username)
}

func (u *userService) GetByID(id string) (models.User, error) {
	return u.Repository.GetSingleUserByID(id)
}

func (u *userService) UpdateUser(user models.User) error {
	return u.Repository.UpdateSingleUser(user)
}

func (u *userService) DeleteUser(id string) (int64, error) {
	return u.Repository.DeleteSingleUser(id)
}
