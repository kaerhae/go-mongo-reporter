package main

import (
	"log"
	"main/reporter-backend/cmd/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	/* GET index */
	router.GET("/", func(ctx *gin.Context) {
		ctx.String(200, "Server up and running!")
	})
	router.GET("/api/reports", routes.GetReports)
	router.POST("/signup", routes.PostNewUser)
	router.POST("/login")

	log.Fatal(router.Run("localhost:3000"))
}
