package routes

import (
	"log"
	"main/cmd/db"
	"main/cmd/middleware"
	"main/cmd/repository"
	"main/cmd/services"
	"main/configs"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	db, _, cancel, err := db.MongoConnect(configs.GetMongoURI())
	if err != nil {
		log.Fatal("Error on init db")
	}
	defer cancel()
	userRepo := repository.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := NewUserHandler(userService)

	reportRepo := repository.NewReportRepository(db)
	reportService := services.NewReportService(reportRepo)
	reportHandler := NewReportRouter(reportService)

	router := gin.Default()

	authorizedGroup := router.Group("/api")
	authorizedGroup.Use(middleware.Authenticate())
	{
		authorizedGroup.GET("/reports", reportHandler.Get)
		authorizedGroup.GET("/reports/:id", reportHandler.GetById)
		authorizedGroup.POST("/reports", reportHandler.Post)
		authorizedGroup.PUT("/reports/:id", reportHandler.Update)
		authorizedGroup.DELETE("/reports/:id", reportHandler.Delete)
	}
	router.GET("/", func(ctx *gin.Context) {
		ctx.String(200, "Server up and running!")
	})

	router.POST("/signup", userHandler.PostNewUser)
	router.POST("/login", userHandler.LoginUser)

	return router
}
