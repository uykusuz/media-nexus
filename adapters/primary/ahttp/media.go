package ahttp

import (
	"context"
	"fmt"
	"media-nexus/adapters/primary/ahttp/ahmodel"
	"media-nexus/httputils"
	"media-nexus/logger"
	"media-nexus/model"
	"media-nexus/services"
	"media-nexus/util"
	"net/http"
)

type mediaEndpoint struct {
	mediaService        services.MediaService
	log                 logger.Logger
	maxUploadFileSizeMB int64
	mediaNameMaxLen     int
	tagIDMaxLen         int
}

//nolint:unused,deadcode
type postMediaRequest struct {
	Name   string   `json:"name"`
	TagIDs []string `json:"tag_ids"`
	// File binary blob
	File []byte `json:"file"`
}

func (e *mediaEndpoint) createContext(r *http.Request) context.Context {
	return util.WithLogger(r.Context(), e.log)
}

// CreateMedia godoc
//
//	@Summary		Create media
//	@Description	create a new media with a list of tags and a name
//	@Tags			media
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			request	body		postMediaRequest	true	"media to be created"
//	@Success		200		{object}	ahmodel.PostMediaResponse
//	@Failure		400		{object}	string
//	@Router			/media [post]
func (e *mediaEndpoint) CreateMedia(w http.ResponseWriter, r *http.Request) {
	ctx := e.createContext(r)

	if err := r.ParseMultipartForm(e.maxUploadFileSizeMB << 20); err != nil {
		http.Error(w, fmt.Sprintf("invalid multipart form: %v", err), http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	if name == "" {
		httputils.RespondWithError(w, http.StatusBadRequest, "File name is required")
		return
	}

	if len(name) > e.mediaNameMaxLen {
		httputils.RespondWithError(w, http.StatusBadRequest, "File name is too long. Maximum is %v", e.mediaNameMaxLen)
		return
	}

	rawTagIDs := r.Form["tag_ids[]"]

	tagIDList := make([]model.TagID, 0, len(rawTagIDs))
	for _, tagID := range rawTagIDs {
		if !e.validateTagID(tagID, w) {
			return
		}
		tagIDList = append(tagIDList, model.TagID(tagID))
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	mediaID, err := e.mediaService.CreateMedia(ctx, name, tagIDList, file)
	if httputils.HandleError(err, w, e.log) {
		return
	}

	response := &ahmodel.PostMediaResponse{
		MediaID: mediaID,
	}

	httputils.RespondWithJSON(http.StatusOK, response, w, e.log, true)
}

func (e *mediaEndpoint) validateTagID(tagID string, w http.ResponseWriter) bool {
	if len(tagID) > e.tagIDMaxLen {
		httputils.RespondWithError(w, http.StatusBadRequest, "Tag ID is too long. Maximum is %v", e.tagIDMaxLen)
		return false
	}

	return true
}

// GetMedia godoc
//
//	@Summary		Query media items
//	@Description	query media items based on some parameters
//	@Tags			media
//	@Produce		json
//	@Param			tag_id	query		string	true	"tag ID to search for"
//	@Success		200		{object}	ahmodel.GetMediaResponse
//	@Failure		400		{object}	string
//	@Router			/media [get]
func (e *mediaEndpoint) GetMedia(w http.ResponseWriter, r *http.Request) {
	ctx := e.createContext(r)

	tagID := r.URL.Query().Get("tag_id")

	if tagID == "" {
		httputils.RespondWithError(w, http.StatusBadRequest, "tag_id parameter is required")
		return
	}

	if !e.validateTagID(tagID, w) {
		return
	}

	mediaItems, err := e.mediaService.FindByTagID(ctx, model.TagID(tagID))
	if httputils.HandleError(err, w, e.log) {
		return
	}

	response := ahmodel.CreateGetMediaResponse(mediaItems)

	httputils.RespondWithJSON(http.StatusOK, response, w, e.log, false)
}
