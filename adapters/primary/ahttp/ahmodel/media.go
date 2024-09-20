package ahmodel

import "media-nexus/model"

type PostMediaResponse struct {
	MediaId string `json:"media_id"`
}

type GetMediaResponse struct {
	Items []*MediaItem
}

type MediaItem struct {
	Name    string   `json:"name"`
	TagIds  []string `json:"tag_ids"`
	FileUrl string   `json:"file_url"`
}

func MediaItemFromModel(item model.MediaItem) *MediaItem {
	return &MediaItem{
		Name:    item.Name(),
		TagIds:  item.TagIds(),
		FileUrl: item.FileUrl(),
	}
}

func CreateGetMediaResponse(items []model.MediaItem) *GetMediaResponse {
	oMediaItems := make([]*MediaItem, 0, len(items))
	for _, mediaItem := range items {
		oMediaItems = append(oMediaItems, MediaItemFromModel(mediaItem))
	}

	return &GetMediaResponse{oMediaItems}
}
