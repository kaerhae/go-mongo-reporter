package routes

import (
	"cmd/reporter-backend/cmd/models"
	"cmd/reporter-backend/cmd/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetReports(c *gin.Context) {
	var reports []*models.Report = service.GetReportsCollection()
	c.IndentedJSON(http.StatusOK, reports)
}
