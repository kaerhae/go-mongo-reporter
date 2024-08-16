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

/*
GET /api/reports route. Allowed all access. Retrieves all reports from db.
*/
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

/*
GET /api/reports route. Allowed all access. Retrieves single report by ID from db.
*/
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

/*
POST /api/reports route. Admin and write permission required. Takes Report model as request body and validates body.

Checks that user exists. If exists, return 400 error.

Finally calls CreateGuestUser method and if successful, returns response with success message.
Also updates reports property from user, which posted new report.
*/
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

/*
PUT /api/reports/:id route. Admin and write permission required. Takes id as url parameter and Report model as request body and validates body.

Checks that report exists. If exists, return 400 error. Checks that user owns report. Return error, if not own nor user is admin.

Finally calls UpdateReport method and if successful, returns response with success message.
*/
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

/*
DELETE /api/reports/:id route. Admin and write permission required. Takes id as url parameter.

Checks that user owns report. Return error, if not own nor user is admin.

Finally calls DeleteReport method and if successful, returns response with success message.
*/
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

/*
Takes a string, and converts it to primitive.ObjectID type
*/
func convertStringToPrimitiveID(id string) (primitive.ObjectID, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	return objID, nil
}

/*
Takes reportRouter type, *gin.Context, and reportID as parameters.

Gets userID and isAdmin boolean from gin.Context. If not exist, return 500 error.

Gets report by given id, if not found, returns 404 error.

Checks that session user is either admin or owns report. If not, return 405 error
*/
func checkOwnership(r *reportRouter, c *gin.Context, reportID string) (int, error) {

	userID, isAdmin, err := middleware.GetSessionData(c)
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
