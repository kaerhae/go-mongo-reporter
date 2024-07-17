package helpers

import (
	"log"
	"main/cmd/models"

	"golang.org/x/crypto/bcrypt"
)

type mockReportRepository struct{}
type mockUserRepository struct{}
type MockUserRepository interface {
	Create(user *models.User) (string, error)
	GetSingleUser(username string) (models.User, error)
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

func (r *mockReportRepository) Create(report *models.Report) error { return nil }

func (r *mockReportRepository) Get() ([]models.Report, error) {
	var list []models.Report
	list = append(list, models.Report{})
	return list, nil
}

func (r *mockReportRepository) GetSingle(id string) (models.Report, error) {
	return models.Report{}, nil
}
func (r *mockReportRepository) Update(newReport *models.Report) error { return nil }
func (r *mockReportRepository) Delete(id string) error                { return nil }
