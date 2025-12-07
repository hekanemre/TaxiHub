package infrastructure

import (
	"context"
	"time"

	"github.com/hekanemre/taxihub/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	DB         *mongo.Database
	Collection string
}

func NewMongoRepository(collection string) (*MongoRepository, error) {
	appConfig := config.Read()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(appConfig.MongoDB.Host))
	if err != nil {
		return nil, err
	}

	db := client.Database(appConfig.MongoDB.DBName)

	return &MongoRepository{
		DB:         db,
		Collection: collection,
	}, nil
}
