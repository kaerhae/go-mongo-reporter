package services

import (
	"log"
	"main/reporter-backend/cmd/db"
	"main/reporter-backend/cmd/models"
	"main/reporter-backend/configs"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetReportsCollection() []*models.Report {
	client, ctx, cancel, err := db.MongoConnect(configs.GetMongoURI())
	if err != nil {
		panic(err)
	}

	var results []*models.Report

	var reports *mongo.Collection = client.Database(configs.GetDBName()).Collection("reports")

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

	defer db.MongoClose(client, ctx, cancel)

	return results

}
