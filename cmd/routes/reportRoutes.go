package routes

import (
	"main/cmd/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReportRouter interface {
	Get(*gin.Context)
	GetById(*gin.Context)
	Post(*gin.Context)
	Update(*gin.Context)
	Delete(*gin.Context)
}

type reportRouter struct {
	Service services.ReportService
}

func NewReportRouter(service services.ReportService) ReportRouter {
	return &reportRouter{Service: service}
}

// Get implements ReportRouter.
func (*reportRouter) Get(*gin.Context) {
	panic("unimplemented")
}

// GetById implements ReportRouter.
func (*reportRouter) GetById(*gin.Context) {
	panic("unimplemented")
}

// Post implements ReportRouter.
func (*reportRouter) Post(*gin.Context) {
	panic("unimplemented")
}

// Update implements ReportRouter.
func (*reportRouter) Update(*gin.Context) {
	panic("unimplemented")
}

// Delete implements ReportRouter.
func (*reportRouter) Delete(*gin.Context) {
	panic("unimplemented")
}

func GetReports(c *gin.Context) {
	reports := services.GetReportsCollection()
	c.IndentedJSON(http.StatusOK, reports)
}
