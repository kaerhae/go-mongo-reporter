package utils

import (
	"errors"
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

func CreateToken(user models.User) (*models.Claims, error) {
	expirationTime := time.Now().Add(5 * time.Minute)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": user.Username,
			"exp":      time.Now().Add(time.Minute * 30).Unix(),
		})

	secretKey := []byte(configs.GetSecret())

	tokenstring, err := token.SignedString(secretKey)
	claims := &models.Claims{
		Username:       user.Username,
		AppRole:        user.AppRole,
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
