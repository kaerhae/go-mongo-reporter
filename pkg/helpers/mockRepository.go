package helpers

import (
	"log"
	"main/pkg/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
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
	Update(newReport *models.Report) error
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

/*
*
*
*  FOR USER TESTS
*
*
 */

type mockUserRepository struct{}

type MockUserRepository interface {
	Create(user *models.User) (string, error)
	GetSingleUser(username string) (models.User, error)
	UpdateUser(newUser models.User) error
	DeleteSingleUser(ID primitive.ObjectID) (int64, error)
}

func InitMockRepository() MockUserRepository {
	return &mockUserRepository{}
}

func (m *mockUserRepository) Create(user *models.User) (string, error) {
	return user.ID.Hex(), nil
}

func (m *mockUserRepository) GetSingleUser(username string) (models.User, error) {
	pwd, err := TestingEnvHashPwd("strong-password")
	if err != nil {
		log.Fatal(err)
	}
	return models.User{
		Username:     username,
		PasswordHash: pwd,
		Email:        "test@test.com",
		AppRole:      "guest",
		CreatedAt:    "2023-01-01",
	}, nil
}

func (m *mockUserRepository) UpdateUser(_ models.User) error                       { return nil }
func (m *mockUserRepository) DeleteSingleUser(_ primitive.ObjectID) (int64, error) { return 1, nil }

func TestingEnvHashPwd(pass string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(pass),
		bcrypt.MinCost,
	)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return string(hash), nil
}
