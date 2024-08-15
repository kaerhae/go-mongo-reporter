package repository

import (
	"context"
	"main/pkg/middleware"
	"main/pkg/models"
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
	Delete(id string) (int64, error)
	UpdateUserReportReferences(userID primitive.ObjectID, reportID primitive.ObjectID) error
}

type reportRepository struct {
	Client *mongo.Database
	Logger middleware.Logger
}

func NewReportRepository(client *mongo.Database, logger middleware.Logger) ReportRepository {
	return &reportRepository{
		Client: client,
		Logger: logger,
	}
}

func (r *reportRepository) Create(report *models.Report) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	collection := r.Client.Collection("reports")
	insertedID, err := collection.InsertOne(ctx, &report)
	if err != nil {
		return "", err
	}
	id := insertedID.InsertedID.(primitive.ObjectID).Hex()
	return id, nil
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

//nolint:dupl
func (r *reportRepository) GetSingle(id string) (models.Report, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	collection := r.Client.Collection("reports")
	var report models.Report
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Report{}, err
	}
	err = collection.FindOne(ctx, bson.D{{
		Key:   "_id",
		Value: objectID,
	}}).Decode(&report)
	if err != nil {
		return models.Report{}, err
	}

	return report, nil
}

// DeleteReport implements ReportRepository.
func (r *reportRepository) Delete(id string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	collection := r.Client.Collection("reports")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return 0, err
	}
	deletedResult, err := collection.DeleteOne(ctx, bson.D{{
		Key:   "_id",
		Value: objectID,
	}})
	if err != nil {
		return 0, err
	}
	return deletedResult.DeletedCount, nil
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

func (r *reportRepository) UpdateUserReportReferences(userID primitive.ObjectID, reportID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	collection := r.Client.Collection("users")
	_, err := collection.UpdateOne(ctx, bson.D{{
		Key:   "_id",
		Value: userID,
	}},
		bson.M{"$addToSet": bson.M{
			"reports": reportID,
		}},
	)

	if err != nil {
		return err
	}

	return nil
}
