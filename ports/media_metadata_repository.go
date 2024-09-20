package ports

import (
	"context"
	"media-nexus/model"
)

type MediaMetadataRepository interface {
	Upsert(ctx context.Context, metadata model.MediaMetadata) error
	Get(ctx context.Context, id model.MediaId) (model.MediaMetadata, error)
	SetUploadComplete(ctx context.Context, metadata model.MediaId, complete bool) error
	FindByTagId(ctx context.Context, id model.TagId) ([]model.MediaMetadata, error)
	FindByChecksum(ctx context.Context, checksum string) (model.MediaMetadata, error)
}
