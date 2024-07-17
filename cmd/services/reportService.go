package services

import (
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

}
