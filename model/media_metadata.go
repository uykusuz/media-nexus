package model

import "time"

type MediaMetadata interface {
	Id() MediaId
	Name() string
	TagIds() []TagId
	Checksum() string
	UploadComplete() bool
	LastUpdate() time.Time
}

func NewMediaMetadata(
	id MediaId,
	name string,
	tagIds []TagId,
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
	id             MediaId
	name           string
	tagIds         []TagId
	checksum       string
	uploadComplete bool
	lastUpdate     time.Time
}

func (m *mediaMetadata) Id() MediaId {
	return m.id
}

func (m *mediaMetadata) Name() string {
	return m.name
}

func (m *mediaMetadata) TagIds() []TagId {
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
