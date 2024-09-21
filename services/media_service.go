package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"io"
	"media-nexus/errortypes"
	"media-nexus/model"
	"media-nexus/ports"
	"media-nexus/util"
	"mime/multipart"
	"time"
)

type MediaService interface {
	CreateMedia(ctx context.Context, name string, tagIds []model.TagId, file multipart.File) (model.MediaId, error)
	FindByTagId(ctx context.Context, tagId model.TagId) ([]model.MediaItem, error)
}

func NewMediaService(
	tags ports.TagRepository,
	mediaMetadata ports.MediaMetadataRepository,
	media ports.MediaRepository,
	mediaUrlLifetime time.Duration,
	incompleteMetadataLifetime time.Duration,
) MediaService {
	return &mediaService{tags, mediaMetadata, media, mediaUrlLifetime, incompleteMetadataLifetime}
}

type mediaService struct {
	tags                       ports.TagRepository
	mediaMetadata              ports.MediaMetadataRepository
	media                      ports.MediaRepository
	mediaUrlLifetime           time.Duration
	incompleteMetadataLifetime time.Duration
}

func (s *mediaService) CreateMedia(
	ctx context.Context,
	name string,
	tagIds []model.TagId,
	file multipart.File,
) (model.MediaId, error) {
	if allExist, err := s.tags.AllExist(ctx, tagIds); err != nil {
		return "", err
	} else if !allExist {
		return "", errortypes.NewBadUserInput("not all tag ids exist. Add them first.")
	}

	metadata, err := createMediaMetadata(name, tagIds, file)
	if err != nil {
		return "", err
	}

	if canProceed, existingMetadataId, err := s.canProceedCreateMedia(ctx, metadata); !canProceed {
		return existingMetadataId, err
	}

	err = s.mediaMetadata.Upsert(ctx, metadata)
	if err != nil {
		return "", err
	}

	err = s.media.CreateMedia(ctx, metadata.Id(), file)
	if err != nil {
		return "", err
	}

	err = s.mediaMetadata.SetUploadComplete(ctx, metadata.Id(), true)
	if err != nil {
		return "", err
	}

	return metadata.Id(), nil
}

func (s *mediaService) canProceedCreateMedia(
	ctx context.Context,
	metadata model.MediaMetadata,
) (bool, model.MediaId, error) {
	existingMetadata, err := s.mediaMetadata.FindByChecksum(ctx, metadata.Checksum())
	if err != nil {
		if !errortypes.IsResourceNotFound(err) {
			return false, "", err
		}

		return true, "", nil
	}

	if existingMetadata.UploadComplete() {
		if metadata.Name() != existingMetadata.Name() {
			return false, existingMetadata.Id(), errortypes.NewResourceAlreadyExistsf(
				"media already exists at id %v but has different name. Expected %v, but is %v",
				existingMetadata.Id(),
				existingMetadata.Name(),
				metadata.Name(),
			)
		}

		if !tagIdsEqual(metadata.TagIds(), existingMetadata.TagIds()) {
			return false, existingMetadata.Id(), errortypes.NewResourceAlreadyExistsf(
				"media already exists at id %v but has different tag Ids",
				existingMetadata.Id(),
			)
		}

		return false, existingMetadata.Id(), nil
	}

	if !existingMetadata.LastUpdate().Add(s.incompleteMetadataLifetime).Before(time.Now()) {
		return false, existingMetadata.Id(), errortypes.NewResourceAlreadyExistsf(
			"media already exists with id '%v'",
			metadata.Id(),
		)
	}

	return true, existingMetadata.Id(), nil
}

func (s *mediaService) FindByTagId(ctx context.Context, tagId model.TagId) ([]model.MediaItem, error) {
	log := util.Logger(ctx)

	metadatas, err := s.mediaMetadata.FindByTagId(ctx, tagId)
	if err != nil {
		return nil, err
	}

	items := make([]model.MediaItem, 0, len(metadatas))

	for _, metadata := range metadatas {
		url, err := s.media.GetMediaUrl(ctx, metadata.Id(), s.mediaUrlLifetime)
		if err != nil {
			log.Errorf("failed to get media url. Adding anyway. Details: %v", err)
		}

		item := model.NewMediaItem(metadata, url)
		items = append(items, item)
	}

	return items, nil
}

func createMediaMetadata(
	name string,
	tagIds []string,
	file multipart.File,
) (model.MediaMetadata, error) {
	checksum, err := computeChecksum(file, sha256.New())
	if err != nil {
		return nil, err
	}

	id := computeHashForMedia(sha256.New(), name, tagIds, checksum)

	return model.NewMediaMetadata(
		id,
		name,
		tagIds,
		checksum,
		false,
		time.Now(),
	), nil
}

func computeChecksum(file multipart.File, hasher hash.Hash) (string, error) {
	defer file.Seek(0, io.SeekStart)

	buffer := make([]byte, 4096)

	for {
		n, err := file.Read(buffer)
		if n > 0 {
			_, err := hasher.Write(buffer[:n])
			if err != nil {
				return "", errortypes.NewInputOutputErrorf("failed to update hash: %v", err)
			}
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return "", errortypes.NewInputOutputErrorf("error reading file: %v", err)
		}
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func computeHashForMedia(hasher hash.Hash, name string, tagIds []string, checksum string) string {
	hasher.Reset()
	hasher.Write([]byte(name))
	for _, tag := range tagIds {
		hasher.Write([]byte(tag))
	}
	hasher.Write([]byte(checksum))

	return hex.EncodeToString(hasher.Sum(nil))
}

func tagIdsEqual(lhs []model.TagId, rhs []model.TagId) bool {
	if len(lhs) != len(rhs) {
		return false
	}

	// yes, of course this in O(n^2), but there should not be tens of thousands of tags
	for _, lhsTag := range lhs {
		foundRhs := false
		for _, rhsTag := range rhs {
			if lhsTag == rhsTag {
				foundRhs = true
			}
		}

		if !foundRhs {
			return false
		}
	}

	return true
}
