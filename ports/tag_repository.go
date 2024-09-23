package ports

import (
	"context"
	"media-nexus/model"
)

type TagRepository interface {
	CreateTag(ctx context.Context, name string) (model.TagID, error)
	ListTags(ctx context.Context) ([]*model.Tag, error)
	DeleteTags(ctx context.Context, ids []model.TagID) error
	AllExist(ctx context.Context, ids []model.TagID) (bool, error)
}
