package ahmodel

import "media-nexus/model"

type PostTagsRequest struct {
	Name string `json:"name"`
}

type PostTagsResponse struct {
	TagID string `json:"tag_id"`
}

type Tag struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func CreateGetTagsResponse(tags []*model.Tag) []*Tag {
	response := make([]*Tag, 0, len(tags))

	for _, tag := range tags {
		oTag := &Tag{tag.ID, tag.Name}
		response = append(response, oTag)
	}

	return response
}
