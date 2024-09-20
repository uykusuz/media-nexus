package amongodb

import (
	"context"
	"media-nexus/errortypes"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoDBClient(ctx context.Context, uri string) (*mongo.Client, error) {
	mongodbClient, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, errortypes.NewUpstreamUnavailablef("mongodb unavailable: %v", err)
	}

	return mongodbClient, nil
}
