package routes

import (
	"log"
	"main/cmd/db"
	"main/cmd/repository"
	"main/cmd/services"
	"main/configs"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	_, db, cancel, err := db.MongoConnect(configs.GetMongoURI())
	if err != nil {
		log.Fatal("Error on init db")
	}
	defer cancel()
	repo := repository.NewUserRepository(db)
	service := services.NewUserService(repo)
	handler := NewUserHandler(service)
	router := gin.Default()

	router.GET("/", func(ctx *gin.Context) {
		ctx.String(200, "Server up and running!")
	})
	router.GET("/api/reports", GetReports)
	router.POST("/signup", handler.PostNewUser)
	router.POST("/login", handler.LoginUser)

	return router
}
