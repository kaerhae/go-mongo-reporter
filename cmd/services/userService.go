package services

import (
	"errors"
	"log"
	"main/cmd/models"
	"main/cmd/repository"
	"main/configs"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(user models.User) (string, error)
	CheckExistingUser(username string) (models.User, error)
	HashPwd(password string) (string, error)
	CheckPassword(hashedPassword string, plainPassword string) bool
	CreateToken(username string, app_role string) (*models.Claims, error)
	DetermineRole(role string) (models.Role, error)
}

type userService struct {
	Repository repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {

	return &userService{Repository: repo}
}

func (u *userService) CreateUser(user models.User) (string, error) {
	r := u.Repository
	addedUser, err := r.Create(&user)
	if err != nil {
		return "", err
	}

	return addedUser, nil
}

func (u *userService) CheckExistingUser(username string) (models.User, error) {
	r := u.Repository
	user, err := r.GetSingleUser(username)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (u *userService) HashPwd(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		bcrypt.MinCost,
	)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return string(hash), nil
}

func (u *userService) CheckPassword(hashedPassword string, plainPassword string) bool {
	bytePlain := []byte(plainPassword)
	byteHashed := []byte(hashedPassword)

	err := bcrypt.CompareHashAndPassword(byteHashed, bytePlain)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func (u *userService) CreateToken(username string, app_role string) (*models.Claims, error) {
	expirationTime := time.Now().Add(5 * time.Minute)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Minute * 30).Unix(),
		})

	secretKey := []byte(configs.GetSecret())

	tokenstring, err := token.SignedString(secretKey)
	claims := &models.Claims{
		Username:       username,
		App_Role:       app_role,
		Token:          tokenstring,
		ExpirationTime: expirationTime,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	if err != nil {
		return nil, err
	}

	return claims, nil
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
