package utils

import (
	"errors"
	"fmt"
	"log"
	"main/configs"
	"main/pkg/models"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func HashPwd(password string) (string, error) {
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

func CheckPassword(hashedPassword string, plainPassword string) error {
	if plainPassword == "" {
		return errors.New("password field was empty")
	}
	bytePlain := []byte(plainPassword)
	byteHashed := []byte(hashedPassword)

	err := bcrypt.CompareHashAndPassword(byteHashed, bytePlain)
	if err != nil {
		return err
	}

	return nil
}

func CreateToken(user models.User) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)

	claims := models.Claims{
		Username:       user.Username,
		AppRole:        user.AppRole,
		ExpirationTime: expirationTime,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		claims)
	secretKey := []byte(configs.GetSecret())
	tokenstring, err := token.SignedString(secretKey)
	fmt.Println("APPROLE CREATE TOKEN", user.AppRole)

	fmt.Println("APPROLE IN CLAIMS: ", claims.AppRole)
	if err != nil {
		return "", err
	}

	return tokenstring, nil
}

func DetermineRole(role string) (models.Role, error) {
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
