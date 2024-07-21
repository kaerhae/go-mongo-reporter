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

func SetupRouter(logger middleware.Logger) *gin.Engine {
	_, db, cancel, err := db.MongoConnect(configs.GetMongoURI())
	if err != nil {
		log.Fatal("Error on init db")
	}
	defer cancel()

	userRepo := repository.NewUserRepository(db, logger)
	userService := services.NewUserService(userRepo, logger)
	userHandler := NewUserHandler(userService, logger)

	reportRepo := repository.NewReportRepository(db, logger)
	reportService := services.NewReportService(reportRepo)
	reportHandler := NewReportRouter(reportService, logger)

	router := gin.Default()

	authorizedGroup := router.Group("/api")
	authorizedGroup.Use(middleware.Authenticate)
	{
		authorizedGroup.GET("/reports", reportHandler.Get)
		authorizedGroup.GET("/reports/:id", reportHandler.GetByID)
		authorizedGroup.POST("/reports", reportHandler.Post)
		authorizedGroup.PUT("/reports/:id", reportHandler.Update)
		authorizedGroup.DELETE("/reports/:id", reportHandler.Delete)
	}
	router.GET("/", func(ctx *gin.Context) {
		logger.LogInfo("Server up and running")
		ctx.String(200, "Server up and running!")
	})

	router.POST("/signup", userHandler.PostNewUser)
	router.POST("/login", userHandler.LoginUser)

	return router
}
