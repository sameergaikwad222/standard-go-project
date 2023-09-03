package database

import (
	"context"
	"log"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongoDB(ctx context.Context) *mongo.Client {
	mongoURI := viper.GetString("MongoURL")
	opts := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatal("Error while connecting Mongo DB", err)
	}
	return client
}
