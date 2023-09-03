package repositories

import (
	"context"
	"errors"
	"fmt"
	"sample-go-crud/app/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SampleRepository struct {
	Collection *mongo.Collection
	Db         *mongo.Database
}

func NewSampleRepository(db *mongo.Database) *SampleRepository {
	collection := db.Collection("sample")
	return &SampleRepository{
		Collection: collection,
		Db:         db,
	}
}

func (s *SampleRepository) GetAllFilteredSamples(ctx context.Context, filter interface{}) ([]*models.Sample, error) {
	if filter == nil {
		filter = bson.M{} // empty query selects all documents in the collection
	}
	opts := options.Find().SetLimit(100)
	cursor, err := s.Collection.Find(ctx, filter, opts)
	if err != nil {
		fmt.Println("Error while getting Samples from mongodb !!")
		panic(err)
	}
	var samples = []*models.Sample{}
	for cursor.Next(ctx) {
		var sample *models.Sample
		err := cursor.Decode(&sample)
		if err != nil {
			return nil, err
		}
		samples = append(samples, sample)
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	cursor.Close(ctx)
	if len(samples) == 0 {
		return nil, errors.New("documents not found")
	}
	return samples, nil
}

func (s *SampleRepository) GetOneSampleById(ctx context.Context, id int) (*models.Sample, error) {
	var sample *models.Sample
	err := s.Collection.FindOne(ctx, bson.M{"id": id}).Decode(&sample)
	if err != nil {
		return nil, err
	}
	return sample, nil
}

func (s *SampleRepository) InsertMultipleSamples(ctx context.Context, samples []interface{}) ([]interface{}, error) {
	if len(samples) == 0 {
		return nil, errors.New("invalid count of docs for insertion")
	}
	result, err := s.Collection.InsertMany(ctx, samples)
	if err != nil {
		return nil, err
	}
	return result.InsertedIDs, nil
}

func (s *SampleRepository) FindOneSampleAndUpdate(ctx context.Context, filter interface{}, update interface{}) (int64, error) {
	result, err := s.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return 0, err
	}
	return result.ModifiedCount, nil
}

func (s *SampleRepository) FindAndUpdateManySamples(ctx context.Context, filter interface{}, update interface{}) (int64, error) {
	result, err := s.Collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return 0, err
	}
	return result.ModifiedCount, nil
}

func (s *SampleRepository) DeleteSingleSampleById(ctx context.Context, id int) (int64, error) {
	result, err := s.Collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}

func (s *SampleRepository) GetFilteredSampleCount(ctx context.Context, filter interface{}) (int64, error) {
	count, err := s.Collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}
	return count, nil
}
