package video

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type MongoRepository struct {
	DB *mongo.Client
}

func NewMongoRepository(client *mongo.Client) *MongoRepository {
	return &MongoRepository{
		DB: client,
	}
}

func (r *MongoRepository) GetAll(ctx context.Context) ([]string, error) {
	return []string{}, nil
}

func (r *MongoRepository) Create(ctx context.Context, ids []string) error {
	return nil
}

func (r *MongoRepository) Delete(ctx context.Context, id string) error {
	return nil
}
