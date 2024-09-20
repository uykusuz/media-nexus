package model

type MediaItem interface {
	MediaMetadata

	FileUrl() string
}

func NewMediaItem(metadata MediaMetadata, fileUrl string) MediaItem {
	return &mediaItem{metadata, fileUrl}
}

type mediaItem struct {
	MediaMetadata
	fileUrl string
}

func (m *mediaItem) FileUrl() string {
	return m.fileUrl
}
