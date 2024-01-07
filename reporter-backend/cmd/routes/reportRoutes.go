package routes

import (
	"main/reporter-backend/cmd/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetReports(c *gin.Context) {
	reports := services.GetReportsCollection()
	c.IndentedJSON(http.StatusOK, reports)
}
