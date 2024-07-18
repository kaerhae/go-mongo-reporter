package services

import (
	"main/cmd/models"
	"main/cmd/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReportService interface {
	GetAllReports() ([]models.Report, error)
	GetSingleReport(id string) (models.Report, error)
	CreateReport(report models.Report) (string, error)
	UpdateReport(newReport models.Report) error
	DeleteReport(id string) (int64, error)
	UpdateReportReferences(userID primitive.ObjectID, reportID primitive.ObjectID) error
}

type reportService struct {
	Repository repository.ReportRepository
}

// CreateReport implements ReportService.
func (r *reportService) CreateReport(report models.Report) (string, error) {
	return r.Repository.Create(&report)
}

// DeleteReport implements ReportService.
func (r *reportService) DeleteReport(id string) (int64, error) {
	return r.Repository.Delete(id)
}

// GetAllReports implements ReportService.
func (r *reportService) GetAllReports() ([]models.Report, error) {
	return r.Repository.Get()
}

// GetSingleReport implements ReportService.
func (r *reportService) GetSingleReport(id string) (models.Report, error) {
	return r.Repository.GetSingle(id)
}

// UpdateReport implements ReportService.
func (r *reportService) UpdateReport(newReport models.Report) error {
	return r.Repository.Update(&newReport)
}

func (r *reportService) UpdateReportReferences(userID primitive.ObjectID, reportID primitive.ObjectID) error {
	return r.Repository.UpdateUserReportReferences(userID, reportID)
}

func NewReportService(repo repository.ReportRepository) ReportService {
	return &reportService{Repository: repo}
}
