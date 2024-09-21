package ports

import (
	"context"
	"io"
	"time"
)

type MediaRepository interface {
	CreateMedia(ctx context.Context, key string, file io.Reader) error
	GetMediaUrl(ctx context.Context, key string, lifetime time.Duration) (string, error)
	DeleteAll(ctx context.Context, keys []string) error
}
