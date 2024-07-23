package middleware

import (
	"errors"
	"main/configs"
	"main/pkg/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func setSecret() ([]byte, error) {
	var secret = configs.GetSecret()
	if secret == "" {
		return nil, errors.New("no secret env set")
	}
	return []byte(secret), nil
}

func Authenticate(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	if auth == "" {
		c.AbortWithStatusJSON(401, gin.H{"message": "401 Unauthorized"})
		return
	}

	secret, err := setSecret()
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"message": "internal server error"})
		return
	}

	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(
		auth,
		claims,
		//nolint:revive
		func(token *jwt.Token) (any, error) {
			return secret, nil
		})
	if err != nil || token == nil {
		c.AbortWithStatusJSON(400, gin.H{"message": "bad request"})
		return
	}

	if !token.Valid {
		c.AbortWithStatusJSON(401, gin.H{"message": "Unauthorized"})
		return
	}

	c.Next()
}
