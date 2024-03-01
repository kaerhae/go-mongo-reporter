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
	Create(report *models.Report) error
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

func (r *reportRepository) Create(report *models.Report) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	collection := r.Client.Collection("reports")
	_, err := collection.InsertOne(ctx, &report)
	if err != nil {
		return err
	}
	return nil
}

func (r *reportRepository) Get() ([]models.Report, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	collection := r.Client.Collection("reports")
	var reports []models.Report

	cur, err := collection.Find(ctx, bson.M{})
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
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Report{}, err
	}
	err = collection.FindOne(ctx, bson.D{{
		Key:   "_id",
		Value: objectId,
	}}).Decode(&report)
	if err != nil {
		return models.Report{}, err
	}

	return report, nil
}

// DeleteReport implements ReportRepository.
func (r *reportRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	collection := r.Client.Collection("reports")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = collection.DeleteOne(ctx, bson.D{{
		Key:   "_id",
		Value: objectId,
	}})
	if err != nil {
		return err
	}
	return nil
}

// UpdateReport implements ReportRepository.
func (r *reportRepository) Update(newReport *models.Report) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	collection := r.Client.Collection("reports")

	_, err := collection.UpdateOne(ctx, bson.D{{
		Key:   "_id",
		Value: newReport.ID,
	}}, bson.M{"$set": newReport})

	if err != nil {
		return err
	}

	return nil
}
