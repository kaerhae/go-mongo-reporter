package repository

import (
	"main/cmd/models"

	"go.mongodb.org/mongo-driver/mongo"
)

type ReportRepository interface {
	Create(report *models.Report) (string, error)
	Get() ([]models.Report, error)
	GetSingle(id string) (models.Report, error)
	Update(newReport *models.Report) error
	Delete(id string) error
}

type reportRepository struct {
	Client *mongo.Database
}

func NewReportRepository(client *mongo.Database) ReportRepository {
	return &reportRepository{Client: client}
}

func (r *reportRepository) Create(report *models.Report) (string, error) {
	panic(report)
}

func (r *reportRepository) Get() ([]models.Report, error) {
	panic("")
}

func (r *reportRepository) GetSingle(id string) (models.Report, error) {
	panic("")
}

// DeleteReport implements ReportRepository.
func (*reportRepository) Delete(id string) error {
	panic("unimplemented")
}

// UpdateReport implements ReportRepository.
func (*reportRepository) Update(newReport *models.Report) error {
	panic("unimplemented")
}
