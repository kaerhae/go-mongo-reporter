package services

import (
	"main/cmd/models"
	"main/cmd/repository"
)

type ReportService interface {
	GetAllReports() ([]models.Report, error)
	GetSingleReport(id string) (models.Report, error)
	CreateReport(report models.Report) (string, error)
	UpdateReport(newReport models.Report) error
	DeleteReport(id string) error
}

type reportService struct {
	Repository repository.ReportRepository
}

// CreateReport implements ReportService.
func (r *reportService) CreateReport(report models.Report) (string, error) {
	return r.Repository.Create(&report)
}

// DeleteReport implements ReportService.
func (*reportService) DeleteReport(id string) error {
	panic("unimplemented")
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
func (*reportService) UpdateReport(newReport models.Report) error {
	panic("unimplemented")
}

func NewReportService(repo repository.ReportRepository) ReportService {
	return &reportService{Repository: repo}
}
