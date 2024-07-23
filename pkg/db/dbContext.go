package db

import (
	"context"
	"fmt"
	"log"
	"main/configs"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func MongoConnect(uri string) (context.Context, *mongo.Database, context.CancelFunc, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("error on connecting db: ", err)
	}
	return ctx, client.Database(configs.GetDBName()), cancel, err
}

func MongoClose(ctx context.Context, client *mongo.Client, cancel context.CancelFunc) {
	defer cancel()

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

func MongoPing(ctx context.Context, client *mongo.Client) error {
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		fmt.Println("Mongo connection error")
		return err
	}

	fmt.Println("Mongo connected succesfully")
	return nil
}
