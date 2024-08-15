package handler

import (
	"errors"
	"fmt"
	"main/pkg/middleware"
	"main/pkg/models"
	"main/pkg/services"

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
	reportID := c.Param("id")
	var body models.Report
	err := c.BindJSON(&body)
	if err != nil {
		r.Logger.LogError(
			fmt.Sprintf("Error happened while binding JSON: %v", err),
		)
		c.IndentedJSON(400, gin.H{"message": "error on parsing request body"})
		return
	}

	existingReport, err := r.Service.GetSingleReport(reportID)
	if err != nil {
		r.Logger.LogInfo("No report found")
		c.IndentedJSON(404, gin.H{"message": "No report found"})
		return
	}

	status, err := checkOwnership(r, c, reportID)
	if err != nil {
		c.IndentedJSON(status, gin.H{"message": "Error while validating request"})
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
	status, err := checkOwnership(r, c, id)
	if err != nil {
		c.IndentedJSON(status, gin.H{"message": fmt.Sprintf("Error while validating request: %v", err)})
		return
	}
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

func getSessionData(c *gin.Context) (id string, isAdmin bool, err error) {

	userID, exists := c.Get("userId")
	if !exists {
		return "", false, errors.New("userId not set")
	}

	id, ok := userID.(string)
	if !ok {
		return "", false, errors.New("userId is not a string")
	}

	isAdminAny, exists := c.Get("isAdmin")
	if !exists {
		return "", false, errors.New("isAdmin not set")
	}

	isAdmin, ok = isAdminAny.(bool)
	if !ok {
		return "", false, errors.New("isAdmin is not a bool")
	}

	return id, isAdmin, nil
}

func checkOwnership(r *reportRouter, c *gin.Context, reportID string) (int, error) {

	userID, isAdmin, err := getSessionData(c)
	if err != nil {
		r.Logger.LogInfo(fmt.Sprintf("Error on fetching session data %v", err))
		return 500, errors.New("internal server error")
	}
	existingReport, err := r.Service.GetSingleReport(reportID)
	if err != nil {
		r.Logger.LogInfo(fmt.Sprintf("No report found, with id: %s", reportID))
		return 404, errors.New("no report found")
	}

	if existingReport.UserID != userID && !isAdmin {
		r.Logger.LogError("Error: requested non-admin user does not own report")
		return 405, errors.New("unauthorized method")
	}

	return -1, nil
}
