package routes

import (
	"main/cmd/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetReports(c *gin.Context) {
	reports, err := services.GetReportsCollection()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
	}
	c.IndentedJSON(http.StatusOK, reports)
}
