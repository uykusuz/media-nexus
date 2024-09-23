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
	CreateMedia(ctx context.Context, name string, tagIDs []model.TagID, file multipart.File) (model.MediaID, error)
	FindByTagID(ctx context.Context, tagID model.TagID) ([]model.MediaItem, error)
}

func NewMediaService(
	tags ports.TagRepository,
	mediaMetadata ports.MediaMetadataRepository,
	media ports.MediaRepository,
	mediaURLLifetime time.Duration,
	incompleteMetadataLifetime time.Duration,
) MediaService {
	return &mediaService{tags, mediaMetadata, media, mediaURLLifetime, incompleteMetadataLifetime}
}

type mediaService struct {
	tags                       ports.TagRepository
	mediaMetadata              ports.MediaMetadataRepository
	media                      ports.MediaRepository
	mediaURLLifetime           time.Duration
	incompleteMetadataLifetime time.Duration
}

func (s *mediaService) CreateMedia(
	ctx context.Context,
	name string,
	tagIds []model.TagID,
	file multipart.File,
) (model.MediaID, error) {
	if allExist, err := s.tags.AllExist(ctx, tagIds); err != nil {
		return "", err
	} else if !allExist {
		return "", errortypes.NewBadUserInput("not all tag ids exist. Add them first.")
	}

	metadata, err := createMediaMetadata(name, tagIds, file)
	if err != nil {
		return "", err
	}

	if canProceed, existingMetadataID, err := s.canProceedCreateMedia(ctx, metadata); !canProceed {
		return existingMetadataID, err
	}

	err = s.mediaMetadata.Upsert(ctx, metadata)
	if err != nil {
		return "", err
	}

	err = s.media.CreateMedia(ctx, metadata.ID(), file)
	if err != nil {
		return "", err
	}

	err = s.mediaMetadata.SetUploadComplete(ctx, metadata.ID(), true)
	if err != nil {
		return "", err
	}

	return metadata.ID(), nil
}

func (s *mediaService) canProceedCreateMedia(
	ctx context.Context,
	metadata model.MediaMetadata,
) (bool, model.MediaID, error) {
	existingMetadata, err := s.mediaMetadata.FindByChecksum(ctx, metadata.Checksum())
	if err != nil {
		if !errortypes.IsResourceNotFound(err) {
			return false, "", err
		}

		return true, "", nil
	}

	if existingMetadata.UploadComplete() {
		if metadata.Name() != existingMetadata.Name() {
			return false, existingMetadata.ID(), errortypes.NewResourceAlreadyExistsf(
				"media already exists at id %v but has different name. Expected %v, but is %v",
				existingMetadata.ID(),
				existingMetadata.Name(),
				metadata.Name(),
			)
		}

		if !tagIdsEqual(metadata.TagIDs(), existingMetadata.TagIDs()) {
			return false, existingMetadata.ID(), errortypes.NewResourceAlreadyExistsf(
				"media already exists at id %v but has different tag Ids",
				existingMetadata.ID(),
			)
		}

		return false, existingMetadata.ID(), nil
	}

	if !existingMetadata.LastUpdate().Add(s.incompleteMetadataLifetime).Before(time.Now()) {
		return false, existingMetadata.ID(), errortypes.NewResourceAlreadyExistsf(
			"media already exists with id '%v'",
			metadata.ID(),
		)
	}

	return true, existingMetadata.ID(), nil
}

func (s *mediaService) FindByTagID(ctx context.Context, tagID model.TagID) ([]model.MediaItem, error) {
	log := util.Logger(ctx)

	metadatas, err := s.mediaMetadata.FindByTagID(ctx, tagID)
	if err != nil {
		return nil, err
	}

	items := make([]model.MediaItem, 0, len(metadatas))

	for _, metadata := range metadatas {
		url, err := s.media.GetMediaURL(ctx, metadata.ID(), s.mediaURLLifetime)
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
	var err error
	defer func() {
		_, sErr := file.Seek(0, io.SeekStart)
		if sErr != nil && err == nil {
			err = sErr
		}
	}()

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

func tagIdsEqual(lhs []model.TagID, rhs []model.TagID) bool {
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
