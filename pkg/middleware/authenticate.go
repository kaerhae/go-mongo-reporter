package middleware

import (
	"fmt"
	"main/configs"
	"main/pkg/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func Authenticate(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	if auth == "" {
		c.AbortWithStatusJSON(401, gin.H{"message": "401 Unauthorized"})
		return
	}

	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(
		auth,
		claims,
		//nolint:revive
		func(token *jwt.Token) (any, error) {
			return []byte(configs.GetSecret()), nil
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

func AuthenticateAdmin(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	if auth == "" {
		c.AbortWithStatusJSON(401, gin.H{"message": "401 Unauthorized"})
		return
	}

	token, err := jwt.ParseWithClaims(
		auth,
		&models.Claims{},
		//nolint:revive
		func(token *jwt.Token) (any, error) {
			return []byte(configs.GetSecret()), nil
		})

	if err != nil || token == nil {
		fmt.Println(err)
		c.AbortWithStatusJSON(400, gin.H{"message": "bad request"})
		return
	}
	claims := token.Claims.(*models.Claims)

	if !token.Valid {
		c.AbortWithStatusJSON(401, gin.H{"message": "Unauthorized"})
		return
	}
	fmt.Printf("CLAIMS IS : %v\n", claims)
	fmt.Printf("APP ROLE IS: %v", claims.AppRole)
	if claims.AppRole != "admin" {
		c.AbortWithStatusJSON(401, gin.H{
			"message": "Permission denied",
		})
	}

	c.Next()
}
