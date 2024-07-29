package api

import (
	"log"
	"main/api/handler"
	"main/configs"
	"main/pkg/db"
	"main/pkg/middleware"
	"main/pkg/repository"
	"main/pkg/services"

	"github.com/gin-gonic/gin"
)

func SetupRouter(logger middleware.Logger) *gin.Engine {
	_, db, cancel, err := db.MongoConnect(configs.GetMongoURI())
	if err != nil {
		log.Fatal("Error on init db", err)
	}
	defer cancel()

	userRepo := repository.NewUserRepository(db)
	userService := services.NewUserService(userRepo, logger)
	userHandler := handler.NewUserHandler(userService, logger)

	reportRepo := repository.NewReportRepository(db, logger)
	reportService := services.NewReportService(reportRepo)
	reportHandler := handler.NewReportRouter(reportService, logger)

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

	adminGroup := router.Group("/user-management")
	adminGroup.Use(middleware.AuthenticateAdmin)
	{
		adminGroup.GET("/users", userHandler.Get)
		adminGroup.GET("/users/:id", userHandler.GetByID)
		adminGroup.POST("/users", userHandler.PostNewUser)
		adminGroup.PUT("/users/:id", userHandler.UpdateUser)
		adminGroup.DELETE("/users/:id", userHandler.DeleteUser)
	}

	router.GET("/", func(ctx *gin.Context) {
		logger.LogInfo("Server up and running")
		ctx.String(200, "Server up and running!")
	})

	router.POST("/login", userHandler.LoginUser)

	return router
}
