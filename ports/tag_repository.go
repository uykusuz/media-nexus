package ports

import (
	"context"
	"media-nexus/model"
)

type TagRepository interface {
	CreateTag(ctx context.Context, name string) (model.TagId, error)
	ListTags(ctx context.Context) ([]*model.Tag, error)
}
