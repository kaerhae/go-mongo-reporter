package helpers

import (
	"main/pkg/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
*
*
*  FOR REPORT TESTS
*
*
 */

type mockReportRepository struct{}

type MockReportRepository interface {
	Create(report *models.Report) (string, error)
	Get() ([]models.Report, error)
	GetSingle(id string) (models.Report, error)
	Delete(id string) (int64, error)
	UpdateUserReportReferences(userID primitive.ObjectID, reportID primitive.ObjectID) error
}

func InitMockReportRepository() MockReportRepository {
	return &mockReportRepository{}
}

func (r *mockReportRepository) Create(report *models.Report) (string, error) {
	return report.ID.Hex(), nil
}

func (r *mockReportRepository) Get() ([]models.Report, error) {
	list := []models.Report{
		{Topic: "test", Author: "test", Description: "test", UserID: primitive.NewObjectID().Hex()},
		{Topic: "test2", Author: "test2", Description: "test2", UserID: primitive.NewObjectID().Hex()},
	}

	return list, nil
}

func (r *mockReportRepository) GetSingle(id string) (models.Report, error) {
	ID, _ := primitive.ObjectIDFromHex(id)
	return models.Report{
		ID:          ID,
		Topic:       "test",
		Author:      "test",
		Description: "test",
		UserID:      primitive.ObjectID{}.Hex(),
	}, nil
}
func (r *mockReportRepository) Update(_ *models.Report) error  { return nil }
func (r *mockReportRepository) Delete(_ string) (int64, error) { return 0, nil }
func (r *mockReportRepository) UpdateUserReportReferences(_ primitive.ObjectID, _ primitive.ObjectID) error {
	return nil
}
