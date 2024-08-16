package middleware

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func GetSessionData(c *gin.Context) (id string, isAdmin bool, err error) {

	userID, exists := c.Get("userId")
	if !exists {
		return "", false, errors.New("userId not set")
	}

	id, ok := userID.(string)
	if !ok {
		return "", false, errors.New("userId is not a string")
	}

	isAdminAny, exists := c.Get("isAdmin")
	if !exists {
		return "", false, errors.New("isAdmin not set")
	}

	isAdmin, ok = isAdminAny.(bool)
	if !ok {
		return "", false, errors.New("isAdmin is not a bool")
	}

	return id, isAdmin, nil
}
