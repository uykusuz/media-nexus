package ammodel

import (
	"media-nexus/errortypes"
	"media-nexus/logger"
	"media-nexus/model"
	"time"
)

type MediaMetadataDocument struct {
	ID             string   `bson:"_id,omitempty"`
	Name           string   `bson:"name,omitempty"`
	TagIDs         []string `bson:"tag_ids,omitempty"`
	Checksum       string   `bson:"checksum,omitempty"`
	UploadComplete bool     `bson:"upload_complete"`
	LastUpdate     string   `bson:"last_update,omitempty"`
}

func NewMediaMetadataDocument(metadata model.MediaMetadata) *MediaMetadataDocument {
	return &MediaMetadataDocument{
		ID:             metadata.ID(),
		Name:           metadata.Name(),
		TagIDs:         metadata.TagIDs(),
		Checksum:       metadata.Checksum(),
		UploadComplete: metadata.UploadComplete(),
		LastUpdate:     LastUpdateToString(metadata.LastUpdate()),
	}
}

func (d *MediaMetadataDocument) ToModel() (model.MediaMetadata, error) {
	t, err := time.Parse(time.RFC3339Nano, d.LastUpdate)
	if err != nil {
		return nil, errortypes.NewInputOutputErrorf("failed to parse last update of media metadata document: %v", err)
	}

	return model.NewMediaMetadata(
		model.MediaID(d.ID),
		d.Name,
		d.TagIDs,
		d.Checksum,
		d.UploadComplete,
		t,
	), nil
}

func MediaMetadataDocumentsToModel(docs []*MediaMetadataDocument, log logger.Logger) ([]model.MediaMetadata, error) {
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

func LastUpdateToString(t time.Time) string {
	return t.Format(time.RFC3339Nano)
}
