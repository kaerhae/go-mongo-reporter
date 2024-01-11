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

	router.GET("/", func(ctx *gin.Context) {
		ctx.String(200, "Server up and running!")
	})

	router.GET("/api/reports", reportHandler.Get)
	router.GET("/api/reports/:id", reportHandler.GetById)
	router.POST("/api/reports", reportHandler.Post)
	router.PUT("/api/reports/:id", reportHandler.Update)
	router.DELETE("/api/reports/:id", reportHandler.Delete)

	router.POST("/signup", userHandler.PostNewUser)
	router.POST("/login", userHandler.LoginUser)

	return router
}
