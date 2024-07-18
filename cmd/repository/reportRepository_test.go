package repository

import (
	"main/cmd/helpers"
	"main/cmd/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreate(t *testing.T) {
	repo := helpers.InitMockReportRepository()
	id := primitive.NewObjectID()
	newReport := models.Report{
		ID:          id,
		Topic:       "Test",
		Description: "Desc",
		Author:      "Tester",
		UserID:      id.Hex(),
	}
	u, err := repo.Create(&newReport)

	if err != nil {
		t.Fail()
	}

	assert.Equal(t, id.Hex(), u)
}

func TestGet(t *testing.T) {
	repo := helpers.InitMockReportRepository()

	u, err := repo.Get()
	if err != nil {
		t.Fail()
	}

	assert.Equal(t, 2, len(u))
}

func TestGetSingle(t *testing.T) {
	repo := helpers.InitMockReportRepository()

	u, err := repo.GetSingle("test")
	if err != nil {
		t.Fail()
	}
	shouldBe := models.Report{
		Topic:       "test",
		Author:      "test",
		Description: "test",
		UserID:      primitive.ObjectID{}.Hex(),
	}
	assert.Equal(t, shouldBe, u)
}

func TestDelete(t *testing.T) {
	repo := helpers.InitMockReportRepository()

	u, err := repo.Delete("test")
	if err != nil {
		t.Fail()
	}

	assert.Equal(t, int64(0), u)
}
