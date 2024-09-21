package ahttp

import (
	"context"
	"fmt"
	"media-nexus/adapters/primary/ahttp/ahmodel"
	"media-nexus/errortypes"
	"media-nexus/httputils"
	"media-nexus/logger"
	"media-nexus/model"
	"media-nexus/services"
	"media-nexus/util"
	"net/http"
	"strings"
)

type mediaEndpoint struct {
	mediaService        services.MediaService
	log                 logger.Logger
	maxUploadFileSizeMB int64
}

func newMediaEndpoint(
	mediaService services.MediaService,
	log logger.Logger,
	maxUploadFileSizeMB int64,
) (*mediaEndpoint, error) {
	if maxUploadFileSizeMB < 1 {
		return nil, errortypes.NewInvalidArgumentf("maxUploadFileSizeMB must be greater than 0.")
	}

	return &mediaEndpoint{mediaService, log, maxUploadFileSizeMB}, nil
}

func (e *mediaEndpoint) HandleMedia(w http.ResponseWriter, r *http.Request) {
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

type postMediaRequest struct {
	Name   string
	TagIds []string
	// File binary blob
	File []byte
}

// handlePost godoc
//
//	@Summary		Create media
//	@Description	create a new media with a list of tags and a name
//	@Tags			media
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			request	body		postMediaRequest	true	"media to be created"
//	@Success		200		{object}	ahmodel.PostMediaResponse
//	@Failure		400		{object}	string
//	@Router			/api/v1/media [post]
func (e *mediaEndpoint) handlePost(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(e.maxUploadFileSizeMB << 20); err != nil {
		http.Error(w, fmt.Sprintf("invalid multipart form: %v", err), http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	if name == "" {
		http.Error(w, "File name is required", http.StatusBadRequest)
		return
	}

	rawTagIds := r.Form["tag_ids[]"]

	tagIdList := make([]model.TagId, 0, len(rawTagIds))
	for _, tagId := range rawTagIds {
		tagIdList = append(tagIdList, model.TagId(tagId))
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	mediaId, err := e.mediaService.CreateMedia(ctx, name, tagIdList, file)
	if httputils.HandleError(err, w, e.log) {
		return
	}

	response := &ahmodel.PostMediaResponse{
		MediaId: mediaId,
	}

	httputils.RespondWithJSON(http.StatusOK, response, w, e.log, true)
}

// handleGet godoc
//
//	@Summary		Query media items
//	@Description	query media items based on some parameters
//	@Tags			media
//	@Produce		json
//	@Param			tag_id	query		string	true	"tag ID to search for"
//	@Success		200		{object}	ahmodel.GetMediaResponse
//	@Failure		400		{object}	string
//	@Router			/api/v1/media [get]
func (e *mediaEndpoint) handleGet(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	tagID := r.URL.Query().Get("tag_id")

	if tagID == "" {
		http.Error(w, "tag_id parameter is required", http.StatusBadRequest)
		return
	}

	mediaItems, err := e.mediaService.FindByTagId(ctx, model.TagId(tagID))
	if httputils.HandleError(err, w, e.log) {
		return
	}

	response := ahmodel.CreateGetMediaResponse(mediaItems)

	httputils.RespondWithJSON(http.StatusOK, response, w, e.log, false)
}
