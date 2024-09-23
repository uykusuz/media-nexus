package ahttp

import (
	"context"
	"media-nexus/adapters/primary/ahttp/ahmodel"
	"media-nexus/httputils"
	"media-nexus/logger"
	"media-nexus/ports"
	"media-nexus/util"
	"net/http"
)

type tagsEndpoint struct {
	tags          ports.TagRepository
	log           logger.Logger
	tagNameMaxLen int
}

func (e *tagsEndpoint) createContext(r *http.Request) context.Context {
	return util.WithLogger(r.Context(), e.log)
}

// CreateTag godoc
//
//	@Summary		Create tag
//	@Description	create a new tag with the given name
//	@Tags			tags
//	@Accept			json
//	@Produce		json
//	@Param			request	body		ahmodel.PostTagsRequest	true	"tag to be created"
//	@Success		200		{object}	ahmodel.PostTagsResponse
//	@Failure		400		{object}	string
//	@Router			/tags [post]
func (e *tagsEndpoint) CreateTag(w http.ResponseWriter, r *http.Request) {
	ctx := e.createContext(r)

	var data ahmodel.PostTagsRequest
	err := httputils.ParseJSONRequestBody(r.Body, &data)
	if httputils.HandleError(err, w, e.log) {
		return
	}

	if len(data.Name) > e.tagNameMaxLen {
		httputils.RespondWithError(w, http.StatusBadRequest, "tag name is too long. Maximum is %v", e.tagNameMaxLen)
		return
	}

	tagId, err := e.tags.CreateTag(ctx, data.Name)
	if httputils.HandleError(err, w, e.log) {
		return
	}

	response := &ahmodel.PostTagsResponse{
		TagId: tagId,
	}

	httputils.RespondWithJSON(http.StatusOK, response, w, e.log, true)
}

// ListTags godoc
//
//	@Summary		List tags
//	@Description	retrieve all tags
//	@Tags			tags
//	@Produce		json
//	@Success		200	{object}	[]ahmodel.Tag
//	@Router			/tags [get]
func (e *tagsEndpoint) ListTags(w http.ResponseWriter, r *http.Request) {
	ctx := e.createContext(r)

	tags, err := e.tags.ListTags(ctx)
	if httputils.HandleError(err, w, e.log) {
		return
	}

	response := ahmodel.CreateGetTagsResponse(tags)

	httputils.RespondWithJSON(http.StatusOK, response, w, e.log, true)
}
