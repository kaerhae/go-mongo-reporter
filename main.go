package main

import (
	"log"
	"main/cmd/db"
	"main/cmd/repository"
	"main/cmd/routes"
	"main/cmd/services"
	"main/configs"

	"github.com/gin-gonic/gin"
)

func main() {

	db, _, cancel, err := db.MongoConnect(configs.GetMongoURI())
	if err != nil {
		log.Fatal("Error on init db")
	}
	defer cancel()
	repo := repository.NewUserRepository(db)
	service := services.NewUserService(repo)
	handler := routes.NewUserHandler(service)
	router := gin.Default()

	/* GET index */
	router.GET("/", func(ctx *gin.Context) {
		ctx.String(200, "Server up and running!")
	})
	router.GET("/api/reports", routes.GetReports)
	router.POST("/signup", handler.PostNewUser)
	router.POST("/login", handler.LoginUser)

	log.Fatal(router.Run("localhost:3000"))
}
