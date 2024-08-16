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

func CreateToken(user models.User) (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)

	claims := models.Claims{
		UserID:      user.ID,
		Username:    user.Username,
		Permissions: user.Permission,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		claims)
	secretKey := []byte(configs.GetSecret())
	tokenstring, err := token.SignedString(secretKey)

	if err != nil {
		return "", err
	}

	return tokenstring, nil
}
