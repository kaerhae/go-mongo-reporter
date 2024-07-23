package routes

import (
	"fmt"
	"main/cmd/middleware"
	"main/cmd/models"
	"main/cmd/services"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReportRouter interface {
	Get(*gin.Context)
	GetByID(*gin.Context)
	Post(*gin.Context)
	Update(*gin.Context)
	Delete(*gin.Context)
}

type reportRouter struct {
	Service services.ReportService
	Logger  middleware.Logger
}

func NewReportRouter(service services.ReportService, logger middleware.Logger) ReportRouter {
	return &reportRouter{
		Service: service,
		Logger:  logger,
	}
}

// Get implements ReportRouter.
func (r *reportRouter) Get(c *gin.Context) {
	reports, err := r.Service.GetAllReports()
	if err != nil {
		r.Logger.LogError(
			fmt.Sprintf("Error happened while fetching reports: %v", err),
		)
		c.IndentedJSON(500, gin.H{"message": fmt.Sprintf("Internal server error: %v", err)})
		return
	}

	c.IndentedJSON(200, reports)
}

// GetById implements ReportRouter.
func (r *reportRouter) GetByID(c *gin.Context) {
	id := c.Param("id")
	report, err := r.Service.GetSingleReport(id)
	if err != nil {
		r.Logger.LogError(
			fmt.Sprintf("Error happened while fetching single report: %v", err),
		)
		c.IndentedJSON(400, gin.H{
			"message": fmt.Sprintf("Error: %v", err),
		})
		return
	}

	c.IndentedJSON(200, report)
}

// Post implements ReportRouter.
func (r *reportRouter) Post(c *gin.Context) {
	var body models.Report
	err := c.BindJSON(&body)
	if err != nil {
		r.Logger.LogError(
			fmt.Sprintf("Error happened while binding JSON: %v", err),
		)
		c.IndentedJSON(500, gin.H{"message": "Internal server error"})
		return
	}

	if body.UserID == "" {
		r.Logger.LogInfo("No user found")
		c.AbortWithStatusJSON(400, gin.H{"message": "No userID found on request"})
		return
	}

	reportID := primitive.NewObjectID()
	userID, err := convertStringToPrimitiveID(body.UserID)
	if err != nil {
		r.Logger.LogError(
			fmt.Sprintf("Error happened while converting ID: %v", err),
		)
		c.AbortWithStatusJSON(500, gin.H{"message": "internal server error"})
		return
	}
	newReport := models.Report{
		ID:          reportID,
		Author:      body.Author,
		Topic:       body.Topic,
		Description: body.Description,
		UserID:      body.UserID,
	}
	_, err = r.Service.CreateReport(newReport)
	if err != nil {
		r.Logger.LogError(
			fmt.Sprintf("Error happened while creating reports: %v", err),
		)
		c.IndentedJSON(400, "error on creating report")
		return
	}

	// Since MongoDB is document db, this function takes care of linking report to user
	err = r.Service.UpdateReportReferences(userID, reportID)
	if err != nil {
		r.Logger.LogError(
			fmt.Sprintf("Error happened while updating report references: %v", err),
		)
		c.IndentedJSON(500, gin.H{"message": "internal server error"})
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
		r.Logger.LogError(
			fmt.Sprintf("Error happened while binding JSON: %v", err),
		)
		c.IndentedJSON(400, gin.H{"message": "error on parsing request body"})
		return
	}

	existingReport, err := r.Service.GetSingleReport(id)
	if err != nil {
		r.Logger.LogInfo("No report found")
		c.IndentedJSON(400, gin.H{"message": "No report found"})
		return
	}

	newReport := models.Report{
		ID:          existingReport.ID,
		Author:      body.Author,
		Topic:       body.Topic,
		Description: body.Description,
		UserID:      existingReport.UserID,
	}
	err = r.Service.UpdateReport(newReport)
	if err != nil {
		r.Logger.LogError(
			fmt.Sprintf("Error happened while updating reports: %v", err),
		)
		c.IndentedJSON(400, gin.H{"message": fmt.Sprintf("Internal server error: %v", err)})
		return
	}

	c.IndentedJSON(200, gin.H{
		"message": fmt.Sprintf("Report \"%s\" was succesfully updated", body.Topic),
	})
}

// Delete implements ReportRouter.
func (r *reportRouter) Delete(c *gin.Context) {
	id := c.Param("id")

	deletedCount, err := r.Service.DeleteReport(id)
	if err != nil {
		r.Logger.LogError(
			fmt.Sprintf("Error happened while deleting reports: %v", err),
		)
		c.IndentedJSON(500, gin.H{"message": "Internal server error"})
		return
	}
	r.Logger.LogInfo(fmt.Sprintf("Deleted %d reports", deletedCount))
	c.IndentedJSON(200, gin.H{
		"message": fmt.Sprintf("Deleted %d reports", deletedCount),
	})
}

func convertStringToPrimitiveID(id string) (primitive.ObjectID, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	return objID, nil
}
