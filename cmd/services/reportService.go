package services

import (
<<<<<<< HEAD
	"main/cmd/models"
	"main/cmd/repository"
)

type ReportService interface {
	GetAllReports() ([]models.Report, error)
	GetSingleReport(id string) (models.Report, error)
	CreateReport(report models.Report) error
	UpdateReport(newReport models.Report) error
	DeleteReport(id string) error
}

type reportService struct {
	Repository repository.ReportRepository
}

// CreateReport implements ReportService.
func (r *reportService) CreateReport(report models.Report) error {
	return r.Repository.Create(&report)
}

// DeleteReport implements ReportService.
func (r *reportService) DeleteReport(id string) error {
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

func NewReportService(repo repository.ReportRepository) ReportService {
	return &reportService{Repository: repo}
=======
	"main/cmd/db"
	"main/cmd/models"
	"main/configs"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetReportsCollection() ([]*models.Report, error) {
	ctx, db, cancel, err := db.MongoConnect(configs.GetMongoURI())
	if err != nil {
		panic(err)
	}

	defer cancel()

	var results []*models.Report

	reports := db.Collection("reports")

	opts := options.Find()

	cur, err := reports.Find(ctx, bson.D{{}}, opts)
	if err != nil {
		return nil, err
	}

	for cur.Next(ctx) {
		var elem models.Report
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}

		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	cur.Close(ctx)

	return results, nil

>>>>>>> master
}
