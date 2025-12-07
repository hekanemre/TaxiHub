package infrastructure

import (
	"context"
	"log"

	"github.com/hekanemre/taxihub/config"
	"github.com/hekanemre/taxihub/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (r *MongoRepository) CreateDriver(ctx context.Context, driver *domain.Driver) error {
	collection := r.DB.Collection(r.Collection)
	_, err := collection.InsertOne(ctx, driver)
	return err
}

func (r *MongoRepository) UpdateDriver(ctx context.Context, driver *domain.Driver) error {
	collection := r.DB.Collection(r.Collection)

	if objID, err := primitive.ObjectIDFromHex(driver.ID); err == nil {
		var driver domain.Driver
		if err := collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&driver); err == nil {
			return err
		}
	}

	filter := bson.M{"_id": driver.ID}
	update := bson.M{"$set": driver}

	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *MongoRepository) GetAllDrivers(ctx context.Context, page, pageSize int) ([]*domain.Driver, error) {
	collection := r.DB.Collection(r.Collection)

	// Calculate skip
	skip := (page - 1) * pageSize

	// Set options for pagination
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(pageSize))

	// Execute query
	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var drivers []*domain.Driver
	for cursor.Next(ctx) {
		var driver domain.Driver
		if err := cursor.Decode(&driver); err != nil {
			return nil, err
		}
		drivers = append(drivers, &driver)
	}

	return drivers, nil
}

func (r *MongoRepository) GetAllDriversNearby(ctx context.Context, lat, lon float64, taxiType string) ([]*domain.Driver, error) {
	appConfig := config.Read()
	maxDistance := appConfig.NearbyDistance // maxDistance in meters

	collection := r.DB.Collection(r.Collection)

	// Check if the 2dsphere index on the location field exists
	indexes, err := collection.Indexes().List(ctx)
	if err != nil {
		return nil, err
	}

	indexExists := false
	for indexes.Next(ctx) {
		var index bson.M
		if err := indexes.Decode(&index); err != nil {
			return nil, err
		}

		// Check if the index is for the 'location' field and of type 2dsphere
		if idxType, ok := index["key"].(bson.M); ok {
			if _, exists := idxType["location"]; exists && index["name"] == "location_2dsphere" {
				indexExists = true
				break
			}
		}
	}

	// If the 2dsphere index doesn't exist, create it
	if !indexExists {
		indexModel := mongo.IndexModel{
			Keys: bson.D{
				{Key: "location", Value: "2dsphere"},
			},
			Options: options.Index().SetName("location_2dsphere"),
		}

		_, err := collection.Indexes().CreateOne(ctx, indexModel)
		if err != nil {
			return nil, err
		}

		log.Println("2dsphere index created successfully on location field")
	}

	filter := bson.M{
		"location": bson.M{
			"$near": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": bson.A{lon, lat},
				},
				"$maxDistance": maxDistance,
			},
		},
		"taxiType": taxiType,
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var drivers []*domain.Driver
	for cursor.Next(ctx) {
		var driver domain.Driver
		if err := cursor.Decode(&driver); err != nil {
			return nil, err
		}
		drivers = append(drivers, &driver)
	}
	return drivers, nil
}

func (r *MongoRepository) GetDriverByID(ctx context.Context, id string) (*domain.Driver, error) {
	collection := r.DB.Collection(r.Collection)

	if objID, err := primitive.ObjectIDFromHex(id); err == nil {
		var driver domain.Driver
		if err := collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&driver); err == nil {
			return &driver, nil
		}
	}

	var driver domain.Driver
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&driver)
	if err != nil {
		return nil, err
	}

	return &driver, nil
}

func (r *MongoRepository) GetDriverByPlate(ctx context.Context, plate string) (*domain.Driver, error) {
	collection := r.DB.Collection(r.Collection)

	var driver domain.Driver
	err := collection.FindOne(ctx, bson.M{"plate": plate}).Decode(&driver)
	if err != nil {
		return nil, err
	}

	return &driver, nil
}
