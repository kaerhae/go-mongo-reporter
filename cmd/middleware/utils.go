package middleware

import (
	"fmt"
	"main/cmd/models"
	"main/configs"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

var secret = []byte(configs.GetSecret())

func Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth := ctx.GetHeader("Authorization")
		fmt.Println(auth)
		if auth == "" {
			ctx.AbortWithStatusJSON(401, gin.H{"message": "401 Unauthorized"})
			return
		}

		claims := &models.Claims{}
		t, err := jwt.ParseWithClaims(auth, claims, func(token *jwt.Token) (any, error) {
			return secret, nil
		})

		if err != nil || t == nil {
			if err == jwt.ErrSignatureInvalid {
				ctx.AbortWithStatusJSON(400, gin.H{"message": "bad request"})
				return
			}
		}

		if !t.Valid {
			ctx.AbortWithStatusJSON(401, gin.H{"message": "Unauthorized"})
			return
		}

	}
}
