package mongostorage

import (
	"context"

	"github.com/ck3g/laterlater/internal/video"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type VideoStorage struct {
	collection *mongo.Collection
}

func NewVideoStorage(client *mongo.Client, dbName string) (*VideoStorage, error) {
	collection := client.Database(dbName).Collection("videos")

	// Create a unique index on the "video_id" field
	indexModel := mongo.IndexModel{
		Keys:    bson.M{"video_id": 1},
		Options: options.Index().SetUnique(true),
	}

	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		return nil, err
	}

	return &VideoStorage{
		collection: collection,
	}, nil
}

func (s *VideoStorage) GetAll(ctx context.Context) ([]video.Video, error) {
	var videos []video.Video

	filter := bson.M{}

	// Find all documents in the collection
	cursor, err := s.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Iterate over the results and decode each document into a Video struct
	for cursor.Next(ctx) {
		var v video.Video
		if err := cursor.Decode(&v); err != nil {
			return nil, err
		}
		videos = append(videos, v)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return videos, nil
}

func (s *VideoStorage) Create(ctx context.Context, videos []video.Video) error {
	for _, v := range videos {
		var iv interface{} = v
		// InsertOne is used instead of InsertMany in order to ignore duplicate errors
		// With InsertMany it's all-or-nothing approach. So if there is 1 out of many documents
		// with the same video_id already in the DB, all the documents are going to be discarded
		_, err := s.collection.InsertOne(ctx, iv)
		if err != nil {
			// Ignore duplicate errors
			if !mongo.IsDuplicateKeyError(err) {
				return err
			}
		}
	}

	return nil
}

func (s *VideoStorage) Delete(ctx context.Context, id string) error {
	videoID := video.ParseID(id)
	filter := bson.D{{Key: "video_id", Value: videoID}}
	_, err := s.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}
