package repository

import (
	"context"
	"main/cmd/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	collection := r.Client.Collection("reports")
	res, err := collection.InsertOne(ctx, &report)
	if err != nil {
		return "", err
	}
	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (r *reportRepository) Get() ([]models.Report, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	collection := r.Client.Collection("reports")
	var reports []models.Report

	cur, err := collection.Find(ctx, "")
	if err != nil {
		return nil, err
	}

	if err = cur.All(context.TODO(), &reports); err != nil {
		panic(err)
	}

	return reports, nil
}
func (r *reportRepository) GetSingle(id string) (models.Report, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	collection := r.Client.Collection("reports")
	var report models.Report

	err := collection.FindOne(ctx, bson.D{{
		Key:   "_id",
		Value: id,
	}}).Decode(&report)
	if err != nil {
		return models.Report{}, nil
	}

	return report, nil
}

// DeleteReport implements ReportRepository.
func (*reportRepository) Delete(id string) error {
	panic("unimplemented")
}

// UpdateReport implements ReportRepository.
func (*reportRepository) Update(newReport *models.Report) error {
	panic("unimplemented")
}
