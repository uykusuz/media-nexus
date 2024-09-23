package amongodb

import (
	"context"
	"media-nexus/adapters/secondary/amongodb/ammodel"
	"media-nexus/errortypes"
	"media-nexus/model"
	"media-nexus/ports"
	"media-nexus/util"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type mediaMetadataRepository struct {
	client                     *mongo.Client
	database                   string
	collection                 string
	incompleteMetadataLifetime time.Duration
}

func NewMediaMetadataRepository(
	client *mongo.Client,
	database string,
	collection string,
	incompleteMetadataLifetime time.Duration,
) (ports.MediaMetadataRepository, util.Runner) {
	repo := &mediaMetadataRepository{client, database, collection, incompleteMetadataLifetime}

	runner := func(ctx context.Context) {
		log := util.Logger(ctx)

		err := repo.ensureIndices(ctx)
		if err != nil {
			log.Errorf("failed to ensure indices for media metadata %v:%v: %v", database, collection, err)
		}
	}

	return repo, runner
}

func (r *mediaMetadataRepository) ensureIndices(ctx context.Context) error {
	collection := r.client.Database(r.database).Collection(r.collection)

	if err := ensureFieldIndex(ctx, collection, "tags_id_index", "tag_ids"); err != nil {
		return err
	}

	if err := ensureFieldIndex(ctx, collection, "checksum_index", "checksum"); err != nil {
		return err
	}

	if err := r.ensureIncompleteMetadataExpireIndex(ctx, collection, "incomplete_metadata_expire_index"); err != nil {
		return err
	}

	return nil
}

func ensureFieldIndex(ctx context.Context, collection *mongo.Collection, indexName string, field string) error {
	background := true

	indexModel := mongo.IndexModel{
		Keys: bson.M{
			field: 1,
		},
		Options: &options.IndexOptions{
			Name:       &indexName,
			Background: &background,
		},
	}

	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	return handleError(err)
}

func (r *mediaMetadataRepository) ensureIncompleteMetadataExpireIndex(
	ctx context.Context,
	collection *mongo.Collection,
	indexName string,
) error {
	// we want to automatically expire metadata zombies. That is metadata docs, for which
	// the upload didn't complete and the last update time is way too old
	ttlIndex := mongo.IndexModel{
		Keys: bson.D{
			{Key: "last_update", Value: 1},
		},
		Options: options.Index().
			SetName(indexName).
			SetExpireAfterSeconds(int32(r.incompleteMetadataLifetime.Seconds())).
			SetPartialFilterExpression(bson.D{
				{Key: "upload_complete", Value: false},
			}),
	}

	_, err := collection.Indexes().CreateOne(ctx, ttlIndex)
	return err
}

func (r *mediaMetadataRepository) Upsert(ctx context.Context, metadata model.MediaMetadata) error {
	doc := ammodel.NewMediaMetadataDocument(metadata)

	// majority writeconcern: this upsert acts as a lock so to say. So we definitely want that written
	collection := r.client.Database(r.database).
		Collection(r.collection, options.Collection().SetWriteConcern(writeconcern.Majority()))

	filter := bson.M{"_id": doc.ID}
	update := bson.M{"$set": doc}
	opts := options.Update().SetUpsert(true)

	_, err := collection.UpdateOne(ctx, filter, update, opts)
	if err := handleError(err); err != nil {
		return err
	}

	return nil
}

func (r *mediaMetadataRepository) SetUploadComplete(ctx context.Context, id model.MediaID, complete bool) error {
	doc := &ammodel.MediaMetadataDocument{
		UploadComplete: complete,
		LastUpdate:     ammodel.LastUpdateToString(time.Now()),
	}

	collection := r.client.Database(r.database).Collection(r.collection)

	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": doc,
	}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err := handleError(err); err != nil {
		return err
	}

	if result.MatchedCount < 1 {
		return errortypes.NewResourceNotFound(id)
	}

	return nil
}

func (r *mediaMetadataRepository) FindByChecksum(ctx context.Context, checksum string) (model.MediaMetadata, error) {
	log := util.Logger(ctx)

	docs, err := r.findDocumentsByChecksum(ctx, checksum)
	if err != nil {
		return nil, err
	}

	if len(docs) < 1 {
		return nil, errortypes.NewResourceNotFoundf("media by checksum %v", checksum)
	}

	if len(docs) > 1 {
		log.Errorf("found multiple documents with same checksum: %v. Will proceed with first one only.", checksum)
	}

	return docs[0].ToModel()
}

func (r *mediaMetadataRepository) FindByTagID(ctx context.Context, id model.TagID) ([]model.MediaMetadata, error) {
	log := util.Logger(ctx)

	docs, err := r.findDocumentsByTagID(ctx, id)
	if err != nil {
		return nil, err
	}

	result := make([]model.MediaMetadata, 0, len(docs))
	for _, doc := range docs {
		metadata, err := doc.ToModel()
		if err != nil {
			log.Errorf("failed to deserialize media metadata: %v", err)
			continue
		}

		result = append(result, metadata)
	}

	return result, nil
}

func (r *mediaMetadataRepository) findDocumentsByChecksum(
	ctx context.Context,
	checksum string,
) ([]*ammodel.MediaMetadataDocument, error) {
	collection := r.client.Database(r.database).Collection(r.collection)

	filter := bson.M{"checksum": checksum}

	cursor, err := collection.Find(ctx, filter)
	if err := handleError(err); err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var docs []*ammodel.MediaMetadataDocument

	for cursor.Next(ctx) {
		var doc ammodel.MediaMetadataDocument
		if err := handleError(cursor.Decode(&doc)); err != nil {
			return nil, err
		}

		docs = append(docs, &doc)
	}

	if err := handleError(cursor.Err()); err != nil {
		return nil, err
	}

	return docs, nil
}

func (r *mediaMetadataRepository) findDocumentsByTagID(
	ctx context.Context,
	id model.TagID,
) ([]*ammodel.MediaMetadataDocument, error) {
	collection := r.client.Database(r.database).Collection(r.collection)

	filter := bson.M{
		"tag_ids": bson.M{
			"$in": []string{id},
		},
	}

	cursor, err := collection.Find(ctx, filter)
	if err := handleError(err); err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var docs []*ammodel.MediaMetadataDocument

	for cursor.Next(ctx) {
		var doc ammodel.MediaMetadataDocument
		if err := handleError(cursor.Decode(&doc)); err != nil {
			return nil, err
		}

		docs = append(docs, &doc)
	}

	if err := handleError(cursor.Err()); err != nil {
		return nil, err
	}

	return docs, nil
}

func (r *mediaMetadataRepository) Get(ctx context.Context, id model.MediaID) (model.MediaMetadata, error) {
	collection := r.client.Database(r.database).Collection(r.collection)

	filter := bson.M{"_id": id}

	var doc ammodel.MediaMetadataDocument

	err := collection.FindOne(ctx, filter).Decode(&doc)
	if err := handleError(err); err != nil {
		return nil, err
	}

	return doc.ToModel()
}

func (r *mediaMetadataRepository) DeleteAll(ctx context.Context, ids []model.MediaID) error {
	collection := r.client.Database(r.database).Collection(r.collection)

	filter := bson.M{"_id": bson.M{"$in": ids}}

	_, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		return errortypes.NewUpstreamCommunicationErrorf("mongodb delete", "failed to delete media metadata: %v", err)
	}

	return err
}
