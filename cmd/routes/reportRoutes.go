package routes

import (
	"main/cmd/services"

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
func (r *reportRouter) Get(c *gin.Context) {
	reports, err := r.Service.GetAllReports()
	if err != nil {
		c.IndentedJSON(500, gin.H{"message": "Internal server error"})
	}

	c.IndentedJSON(200, reports)
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
