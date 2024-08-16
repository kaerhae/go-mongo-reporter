package middleware

import (
	"main/configs"
	"main/pkg/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

/*
Middleware function, which authorizes admins and non-admin users.
Function gets token from Authorization header, checks validity, and finally checks
if user has permission to perform operation
*/
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

	/*
		If user has no admin admin rights, permissions are checked by method
	*/
	if !claims.Permissions.Admin {
		switch c.Request.Method {
		case http.MethodGet:
			if !claims.Permissions.Read {
				c.AbortWithStatusJSON(401, gin.H{"message": "401 Unauthorized"})
				return
			}
		case http.MethodPost:
			if !claims.Permissions.Write {
				c.AbortWithStatusJSON(401, gin.H{"message": "401 Unauthorized"})
				return
			}
		case http.MethodPut:
			if !claims.Permissions.Write {
				c.AbortWithStatusJSON(401, gin.H{"message": "401 Unauthorized"})
				return
			}
		case http.MethodDelete:
			if !claims.Permissions.Write {
				c.AbortWithStatusJSON(401, gin.H{"message": "401 Unauthorized"})
				return
			}
		default:
			c.AbortWithStatusJSON(500, "internal server error")
		}
	}
	/*
		Finally when authorization is successful, add userId to store.
		If admin, set isAdmin == true to store.
	*/
	c.Set("userId", claims.UserID.Hex())
	if claims.Permissions.Admin {
		c.Set("isAdmin", true)
	}
	c.Next()
}

func AuthenticateTokenOnly(c *gin.Context) {
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

	/*
		Finally when authorization is successful, add userId to store.
		If admin, set isAdmin == true to store.
	*/
	c.Set("userId", claims.UserID.Hex())
	if claims.Permissions.Admin {
		c.Set("isAdmin", true)
	}
	c.Next()
}

/*
Middleware function authorizes admin-only.
Gets token from Authorization header, validates it, and checks that user is admin.
*/
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
		c.AbortWithStatusJSON(400, gin.H{"message": "bad request"})
		return
	}

	if !token.Valid {
		c.AbortWithStatusJSON(401, gin.H{"message": "Unauthorized"})
		return
	}
	claims := token.Claims.(*models.Claims) //nolint:errcheck
	if claims == nil {
		c.AbortWithStatusJSON(401, gin.H{"message": "Unauthorized"})
		return
	}

	if !claims.Permissions.Admin {
		c.AbortWithStatusJSON(403, gin.H{
			"message": "Permission denied",
		})
	}

	/* Finally when authorization is successful, add userId and isAdmin == true to store */
	c.Set("userId", claims.UserID.Hex())
	c.Set("isAdmin", true)
	c.Next()
}
