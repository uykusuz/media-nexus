package amongodb

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"media-nexus/errortypes"
	"media-nexus/model"
	"media-nexus/ports"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type tagDocument struct {
	ID   string `bson:"_id"`
	Name string `bson:"name"`
}

func NewTagRepository(client *mongo.Client, database string, collection string) ports.TagRepository {
	return &tagRepository{client, database, collection}
}

type tagRepository struct {
	client     *mongo.Client
	database   string
	collection string
}

func (r *tagRepository) CreateTag(ctx context.Context, name string) (model.TagID, error) {
	id, err := r.insertTagIfNotExists(ctx, name)
	if err != nil {
		return "", err
	}

	return id, nil
}

func createIDForName(name string) (string, error) {
	hasher := sha256.New()
	if _, err := hasher.Write([]byte(name)); err != nil {
		return "", errortypes.NewInputOutputErrorf("failed to hash %v", name)
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func (r *tagRepository) insertTagIfNotExists(ctx context.Context, name string) (string, error) {
	collection := r.client.Database(r.database).Collection(r.collection)

	filter := bson.M{"name": name}

	newID, err := createIDForName(name)
	if err != nil {
		return "", err
	}

	newTagDoc := &tagDocument{
		ID:   newID,
		Name: name,
	}

	update := bson.M{
		"$setOnInsert": newTagDoc,
	}

	opts := options.Update().SetUpsert(true)

	result, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return "", errortypes.NewUpstreamCommunicationErrorf("mongodb", "failed to upsert tag: %v", err)
	}

	if result.UpsertedCount > 0 {
		if oid, ok := result.UpsertedID.(string); ok {
			return oid, nil
		}

		return "", errortypes.NewUpstreamCommunicationErrorf(
			"mongodb",
			"failed to assert UpsertedID as string",
		)
	}

	// if no document was inserted, find the existing document by name

	var existingTagDoc tagDocument
	err = collection.FindOne(ctx, filter).Decode(&existingTagDoc)
	if err != nil {
		return "", errortypes.NewUpstreamCommunicationErrorf(
			"mongodb",
			"failed to find existing tag: %v",
			err,
		)
	}

	return existingTagDoc.ID, nil
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
		var tag tagDocument
		err := cursor.Decode(&tag)
		if err != nil {
			return nil, errortypes.NewInputOutputErrorf("failed to decode mongodb tag: %v", err)
		}

		mTag := &model.Tag{
			ID:   tag.ID,
			Name: tag.Name,
		}

		tags = append(tags, mTag)
	}

	if err := cursor.Err(); err != nil {
		return nil, errortypes.NewInputOutputErrorf("error during cursor iteration: %v", err)
	}

	return tags, nil
}

func (r *tagRepository) DeleteTags(ctx context.Context, tagIds []model.TagID) error {
	collection := r.client.Database(r.database).Collection(r.collection)

	filter := bson.M{"_id": bson.M{"$in": tagIds}}

	_, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return errortypes.NewUpstreamCommunicationErrorf("mongodb delete", "failed to delete tags '%v': %v", tagIds, err)
	}

	return err
}

func (r *tagRepository) AllExist(ctx context.Context, ids []model.TagID) (bool, error) {
	collection := r.client.Database(r.database).Collection(r.collection)

	filter := bson.M{"_id": bson.M{"$in": ids}}

	count, err := collection.CountDocuments(ctx, filter)
	if err := handleError(err); err != nil {
		return false, err
	}

	return count == int64(len(ids)), nil
}
