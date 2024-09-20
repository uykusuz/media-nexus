package ammodel

import (
	"media-nexus/errortypes"
	"media-nexus/logger"
	"media-nexus/model"
	"time"
)

type MediaMetadataDocument struct {
	Id             string   `bson:"_id,omitempty"`
	Name           string   `bson:"name,omitempty"`
	TagIds         []string `bson:"tag_ids,omitempty"`
	Checksum       string   `bson:"checksum,omitempty"`
	UploadComplete bool     `bson:"upload_complete"`
	LastUpdate     string   `bson:"last_update,omitempty"`
}

func NewMediaMetadataDocument(metadata model.MediaMetadata) *MediaMetadataDocument {
	return &MediaMetadataDocument{
		Id:             metadata.Id(),
		Name:           metadata.Name(),
		TagIds:         metadata.TagIds(),
		Checksum:       metadata.Checksum(),
		UploadComplete: metadata.UploadComplete(),
		LastUpdate:     metadata.LastUpdate().Format(time.RFC3339Nano),
	}
}

func (d *MediaMetadataDocument) ToModel() (model.MediaMetadata, error) {
	t, err := time.Parse(time.RFC3339Nano, d.LastUpdate)
	if err != nil {
		return nil, errortypes.NewInputOutputErrorf("failed to parse last update of media metadata document: %v", err)
	}

	return model.NewMediaMetadata(
		model.MediaId(d.Id),
		d.Name,
		d.TagIds,
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
