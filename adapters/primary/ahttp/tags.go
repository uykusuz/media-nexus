package ahttp

import (
	"context"
	"media-nexus/adapters/primary/ahttp/ahmodel"
	"media-nexus/httputils"
	"media-nexus/logger"
	"media-nexus/ports"
	"media-nexus/util"
	"net/http"
	"strings"
)

type tagsEndpoint struct {
	tags ports.TagRepository
	log  logger.Logger
}

func newTagsEndpoint(
	tags ports.TagRepository,
	log logger.Logger,
) (*tagsEndpoint, error) {
	return &tagsEndpoint{tags, log}, nil
}

func (e *tagsEndpoint) HandleTags(w http.ResponseWriter, r *http.Request) {
	ctx := util.WithLogger(r.Context(), e.log)

	switch r.Method {
	case http.MethodPost:
		e.handlePost(ctx, w, r)
	case http.MethodGet:
		e.handleGet(ctx, w, r)
	default:
		w.Header().Set("Allow", strings.Join([]string{http.MethodPost, http.MethodGet}, ","))
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
}

// handlePost godoc
//
//	@Summary		Create tag
//	@Description	create a new tag with the given name
//	@Tags			tags
//	@Accept			json
//	@Produce		json
//	@Param			request	body		ahmodel.PostTagsRequest	true	"tag to be created"
//	@Success		200		{object}	ahmodel.PostTagsResponse
//	@Failure		400		{object}	string
//	@Router			/api/v1/tags [post]
func (e *tagsEndpoint) handlePost(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var data ahmodel.PostTagsRequest
	err := httputils.ParseJSONRequestBody(r.Body, &data)
	if httputils.HandleError(err, w, e.log) {
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

// handleGet godoc
//
//	@Summary		List tags
//	@Description	retrieve all tags
//	@Tags			tags
//	@Produce		json
//	@Success		200	{object}	[]ahmodel.Tag
//	@Router			/api/v1/tags [get]
func (e *tagsEndpoint) handleGet(ctx context.Context, w http.ResponseWriter, _ *http.Request) {
	tags, err := e.tags.ListTags(ctx)
	if httputils.HandleError(err, w, e.log) {
		return
	}

	response := ahmodel.CreateGetTagsResponse(tags)

	httputils.RespondWithJSON(http.StatusOK, response, w, e.log, true)
}
