package model

import "time"

type MediaMetadata interface {
	ID() MediaID
	Name() string
	TagIDs() []TagID
	Checksum() string
	UploadComplete() bool
	LastUpdate() time.Time
}

func NewMediaMetadata(
	id MediaID,
	name string,
	tagIds []TagID,
	checksum string,
	uploadComplete bool,
	lastUpdate time.Time,
) MediaMetadata {
	return &mediaMetadata{
		id:             id,
		name:           name,
		tagIds:         tagIds,
		checksum:       checksum,
		uploadComplete: uploadComplete,
		lastUpdate:     lastUpdate,
	}
}

type mediaMetadata struct {
	id             MediaID
	name           string
	tagIds         []TagID
	checksum       string
	uploadComplete bool
	lastUpdate     time.Time
}

func (m *mediaMetadata) ID() MediaID {
	return m.id
}

func (m *mediaMetadata) Name() string {
	return m.name
}

func (m *mediaMetadata) TagIDs() []TagID {
	return m.tagIds
}

func (m *mediaMetadata) Checksum() string {
	return m.checksum
}

func (m *mediaMetadata) UploadComplete() bool {
	return m.uploadComplete
}

func (m *mediaMetadata) LastUpdate() time.Time {
	return m.lastUpdate
}
