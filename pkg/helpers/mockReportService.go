package helpers

import (
	"main/pkg/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MockReportService struct {
	Repository MockReportRepository
}

func (s *MockReportService) GetAllReports() ([]models.Report, error) {
	return []models.Report{
		{Topic: "test", Author: "Test"},
	}, nil
}

func (s *MockReportService) GetSingleReport(id string) (models.Report, error) {
	objID, _ := primitive.ObjectIDFromHex("123456789012345678901234")
	return models.Report{
		ID:     objID,
		Topic:  "test",
		Author: "Test",
		UserID: "123456789012345678901234",
	}, nil
}

func (s *MockReportService) CreateReport(_ models.Report) (string, error) {
	return "", nil
}

func (s *MockReportService) UpdateReport(newReport models.Report) error {
	return nil
}

func (s *MockReportService) DeleteReport(id string) (int64, error) {
	return int64(0), nil
}

func (s *MockReportService) UpdateReportReferences(userID primitive.ObjectID, _ primitive.ObjectID) error {
	return nil
}
