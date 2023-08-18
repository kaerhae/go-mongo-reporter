package main

import (
	"cmd/reporter-backend/cmd/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	/* GET index */
	router.GET("/reports", routes.GetReports)
	router.POST("/users", routes.PostNewUser)

	router.Run("localhost:3000")
}
