package amongodb

import (
	"context"
	"media-nexus/errortypes"
	"media-nexus/model"
	"media-nexus/ports"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type tag struct {
	ID   primitive.ObjectID `bson:"_id"  json:"id"`   // MongoDB ObjectId
	Name string             `bson:"name" json:"name"` // Tag name
}

func NewTagRepository(client *mongo.Client, database string, collection string) ports.TagRepository {
	return &tagRepository{client, database, collection}
}

type tagRepository struct {
	client     *mongo.Client
	database   string
	collection string
}

func (r *tagRepository) CreateTag(ctx context.Context, name string) (model.TagId, error) {
	id, err := r.insertTagIfNotExists(ctx, name)
	if err != nil {
		return "", err
	}

	return id.Hex(), nil
}

func (r *tagRepository) insertTagIfNotExists(ctx context.Context, name string) (primitive.ObjectID, error) {
	collection := r.client.Database(r.database).Collection(r.collection)

	filter := bson.M{"name": name}

	update := bson.M{
		"$setOnInsert": bson.M{"name": name},
	}

	opts := options.Update().SetUpsert(true)

	result, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return primitive.NilObjectID, errortypes.NewUpstreamCommunicationErrorf("mongodb", "failed to upsert tag: %v", err)
	}

	if result.UpsertedCount < 1 {
		if oid, ok := result.UpsertedID.(primitive.ObjectID); ok {
			return oid, nil
		}

		return primitive.NilObjectID, errortypes.NewUpstreamCommunicationErrorf(
			"mongodb",
			"failed to assert UpsertedID as ObjectID",
		)
	}

	// If no document was inserted, find the existing document by name

	var existingDocument struct {
		ID primitive.ObjectID `bson:"_id"`
	}

	err = collection.FindOne(ctx, filter).Decode(&existingDocument)
	if err != nil {
		return primitive.NilObjectID, errortypes.NewUpstreamCommunicationErrorf(
			"mongodb",
			"failed to find existing tag: %v",
			err,
		)
	}

	return existingDocument.ID, nil
}

func (r *tagRepository) ListTags(ctx context.Context) ([]*model.Tag, error) {
	collection := r.client.Database(r.database).Collection(r.collection)

	filter := bson.M{}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, errortypes.NewUpstreamCommunicationErrorf("mongodb find", "failed to find tags: %v", err)
	}

	defer cursor.Close(ctx)

	var tags []*model.Tag

	for cursor.Next(ctx) {
		var tag tag
		err := cursor.Decode(&tag)
		if err != nil {
			return nil, errortypes.NewInputOutputErrorf("failed to decode mongodb tag: %v", err)
		}

		mTag := &model.Tag{
			Id:   tag.ID.Hex(),
			Name: tag.Name,
		}

		tags = append(tags, mTag)
	}

	if err := cursor.Err(); err != nil {
		return nil, errortypes.NewInputOutputErrorf("error during cursor iteration: %v", err)
	}

	return tags, nil
}
