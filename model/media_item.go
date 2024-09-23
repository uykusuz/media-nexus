package model

type MediaItem interface {
	MediaMetadata

	FileURL() string
}

func NewMediaItem(metadata MediaMetadata, fileURL string) MediaItem {
	return &mediaItem{metadata, fileURL}
}

type mediaItem struct {
	MediaMetadata
	fileURL string
}

func (m *mediaItem) FileURL() string {
	return m.fileURL
}
