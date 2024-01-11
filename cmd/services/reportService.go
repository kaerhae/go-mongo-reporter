package services

import (
	"log"
	"main/cmd/db"
	"main/cmd/models"
	"main/cmd/repository"
	"main/configs"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ReportService interface {
	GetAllReports() ([]models.Report, error)
	GetSingleReport(id string) (models.Report, error)
	CreateReport(report models.Report) (models.Report, error)
	UpdateReport(newReport models.Report) error
	DeleteReport(id string) error
}

type reportService struct {
	Repository repository.ReportRepository
}

// CreateReport implements ReportService.
func (*reportService) CreateReport(report models.Report) (models.Report, error) {
	panic("unimplemented")
}

// DeleteReport implements ReportService.
func (*reportService) DeleteReport(id string) error {
	panic("unimplemented")
}

// GetAllReports implements ReportService.
func (*reportService) GetAllReports() ([]models.Report, error) {
	panic("unimplemented")
}

// GetSingleReport implements ReportService.
func (*reportService) GetSingleReport(id string) (models.Report, error) {
	panic("unimplemented")
}

// UpdateReport implements ReportService.
func (*reportService) UpdateReport(newReport models.Report) error {
	panic("unimplemented")
}

func NewReportService(repo repository.ReportRepository) ReportService {
	return &reportService{Repository: repo}
}

func GetReportsCollection() []*models.Report {
	db, ctx, cancel, err := db.MongoConnect(configs.GetMongoURI())
	if err != nil {
		panic(err)
	}

	defer cancel()

	var results []*models.Report

	var reports *mongo.Collection = db.Collection("reports")

	opts := options.Find()

	cur, err := reports.Find(ctx, bson.D{{}}, opts)
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(ctx) {
		var elem models.Report
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}

		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(ctx)

	return results

}
