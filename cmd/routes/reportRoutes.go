package routes

import (
	"fmt"
	"main/cmd/models"
	"main/cmd/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

<<<<<<< HEAD
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
		c.IndentedJSON(500, gin.H{"message": fmt.Sprintf("Internal server error: %v", err)})
		return
	}

	c.IndentedJSON(200, reports)
}

// GetById implements ReportRouter.
func (r *reportRouter) GetById(c *gin.Context) {
	id := c.Param("id")
	report, err := r.Service.GetSingleReport(id)
	if err != nil {
		c.IndentedJSON(400, gin.H{"message": "Error"})
		return
	}

	c.IndentedJSON(200, report)
}

// Post implements ReportRouter.
func (r *reportRouter) Post(c *gin.Context) {
	var body models.Report
	err := c.BindJSON(&body)
	if err != nil {
		c.IndentedJSON(500, gin.H{"message": "Internal server error"})
		return
	}

	newReport := models.Report{
		ID:          primitive.NewObjectID(),
		Author:      body.Author,
		Topic:       body.Topic,
		Description: body.Description,
	}
	err = r.Service.CreateReport(newReport)
	if err != nil {
		c.IndentedJSON(400, "err")
		return
	}

	c.IndentedJSON(200, gin.H{"message": "Report was succesfully created"})
}

// Update implements ReportRouter.
func (r *reportRouter) Update(c *gin.Context) {
	id := c.Param("id")
	var body models.Report
	err := c.BindJSON(&body)
	if err != nil {
		c.IndentedJSON(400, gin.H{"message": "Internal server error"})
		return
	}

	existingReport, err := r.Service.GetSingleReport(id)
	if err != nil {
		c.IndentedJSON(400, gin.H{"message": "No user found"})
		return
	}

	newReport := models.Report{
		ID:          existingReport.ID,
		Author:      body.Author,
		Topic:       body.Topic,
		Description: body.Description,
	}
	err = r.Service.UpdateReport(newReport)
	if err != nil {
		c.IndentedJSON(400, gin.H{"message": fmt.Sprintf("Internal server error: %v", err)})
		return
	}

	c.IndentedJSON(200, fmt.Sprintf("Report \"%s\" was succesfully updated", body.Topic))
}

// Delete implements ReportRouter.
func (r *reportRouter) Delete(c *gin.Context) {
	id := c.Param("id")

	err := r.Service.DeleteReport(id)
	if err != nil {
		c.IndentedJSON(500, gin.H{"message": "Internal server error"})
		return
	}

	c.IndentedJSON(200, fmt.Sprintf("Report \"%s\" was succesfully deleted", id))
=======
func GetReports(c *gin.Context) {
	reports, err := services.GetReportsCollection()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, err)
	}
	c.IndentedJSON(http.StatusOK, reports)
>>>>>>> master
}
