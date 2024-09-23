package ports

import (
	"context"
	"media-nexus/model"
)

type MediaMetadataRepository interface {
	Upsert(ctx context.Context, metadata model.MediaMetadata) error
	Get(ctx context.Context, id model.MediaID) (model.MediaMetadata, error)
	SetUploadComplete(ctx context.Context, metadata model.MediaID, complete bool) error
	FindByTagID(ctx context.Context, id model.TagID) ([]model.MediaMetadata, error)
	FindByChecksum(ctx context.Context, checksum string) (model.MediaMetadata, error)
	DeleteAll(ctx context.Context, ids []model.MediaID) error
}
