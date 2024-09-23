package ihttp

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"media-nexus/adapters/primary/ahttp/ahmodel"
	"media-nexus/integrationtests"
	"media-nexus/model"
	"mime/multipart"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type mediaE2ETestSuite struct {
	integrationtests.E2ETestSuite
}

func TestMedia(t *testing.T) {
	suite.Run(t, &mediaE2ETestSuite{})
}

func (s *mediaE2ETestSuite) TestCreateMedia() {
	ctx := s.Context()

	tagIDs := s.createTags(ctx, 2)
	defer func() { s.LogIfError(s.App().TagRepo().DeleteTags(ctx, tagIDs), "delete tags") }()

	mediaID := s.createMedia(s.GenerateAlphanumeric(10), tagIDs, "./../assets/test.png")

	defer func() {
		s.LogIfError(s.App().MediaMetadataRepo().DeleteAll(ctx, []model.MediaID{mediaID}), "delete all media metadata")
	}()
	defer func() { s.LogIfError(s.App().MediaRepo().DeleteAll(ctx, []string{mediaID}), "delete all medias") }()
}

func (s *mediaE2ETestSuite) TestCreateMediaIdempotency() {
	ctx := s.Context()

	tagIDs := s.createTags(ctx, 2)
	defer func() { s.LogIfError(s.App().TagRepo().DeleteTags(ctx, tagIDs), "delete tags") }()

	name := s.GenerateAlphanumeric(10)
	mediaID := s.createMedia(name, tagIDs, "./../assets/test.png")
	defer func() {
		s.LogIfError(s.App().MediaMetadataRepo().DeleteAll(ctx, []model.MediaID{mediaID}), "delete media metadata")
	}()
	defer func() { s.LogIfError(s.App().MediaRepo().DeleteAll(ctx, []string{mediaID}), "delete media") }()

	mediaID2 := s.createMedia(name, tagIDs, "./../assets/test.png")

	if !s.Equal(mediaID, mediaID2) {
		defer func() {
			s.LogIfError(s.App().MediaMetadataRepo().DeleteAll(ctx, []model.MediaID{mediaID2}), "delete all media metadata")
		}()
		defer func() { s.LogIfError(s.App().MediaRepo().DeleteAll(ctx, []string{mediaID2}), "delete all media") }()
	}
}

func (s *mediaE2ETestSuite) TestGetMediaByTagIds() {
	ctx := s.Context()

	tagIDs := s.createTags(ctx, 2)
	defer func() { s.LogIfError(s.App().TagRepo().DeleteTags(ctx, tagIDs), "delete tags") }()

	mediaID := s.createMedia(s.GenerateAlphanumeric(10), tagIDs, "./../assets/test.png")
	mediaID2 := s.createMedia(s.GenerateAlphanumeric(10), []model.TagID{tagIDs[0]}, "./../assets/test2.png")

	mediaIDs := []model.MediaID{mediaID, mediaID2}

	defer func() { s.LogIfError(s.App().MediaMetadataRepo().DeleteAll(ctx, mediaIDs), "delete media metadatas") }()
	defer func() { s.LogIfError(s.App().MediaRepo().DeleteAll(ctx, mediaIDs), "delete medias") }()

	mediaItems := s.getMedia(tagIDs[0])
	s.Equal(2, len(mediaItems))

	mediaItems = s.getMedia(tagIDs[1])
	s.Equal(1, len(mediaItems))
}

func (s *mediaE2ETestSuite) createTags(ctx context.Context, count int) []model.TagID {
	var tagIDs []model.TagID

	for i := 0; i < count; i++ {
		tagID, err := s.App().TagRepo().CreateTag(ctx, s.GenerateAlphanumeric(10))
		s.Require().NoError(err)
		tagIDs = append(tagIDs, tagID)
	}

	return tagIDs
}

func (s *mediaE2ETestSuite) createMedia(name string, tagIds []model.TagID, filePath string) model.MediaID {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	err := writer.WriteField("name", name)
	s.Require().NoError(err)

	for _, tagID := range tagIds {
		err := writer.WriteField("tag_ids[]", tagID)
		s.Require().NoError(err)
	}

	file, err := os.Open(filePath)
	s.Require().NoError(err)
	defer file.Close()

	fileWriter, err := writer.CreateFormFile("file", filePath)
	s.Require().NoError(err)

	_, err = io.Copy(fileWriter, file)
	s.Require().NoError(err)

	err = writer.Close()
	s.Require().NoError(err)

	req, err := http.NewRequest("POST", s.CreateServerURL("/media"), &body)
	s.Require().NoError(err)

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := s.Client().Do(req)
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusOK, resp.StatusCode)

	var postMediaResponse ahmodel.PostMediaResponse
	decoder := json.NewDecoder(resp.Body)
	s.NoError(decoder.Decode(&postMediaResponse))

	s.Require().NotEmpty(postMediaResponse.MediaID)

	return postMediaResponse.MediaID
}

func (s *mediaE2ETestSuite) getMedia(tagID model.TagID) []*ahmodel.MediaItem {
	url := s.CreateServerURL("/media?tag_id=%v", tagID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	s.NoError(err)

	req.Header.Set("Content-Type", "application/json")

	response, err := s.Client().Do(req)
	s.NoError(err)

	defer response.Body.Close()

	s.Equal(http.StatusOK, response.StatusCode)

	var getMediaResponse ahmodel.GetMediaResponse
	decoder := json.NewDecoder(response.Body)
	s.NoError(decoder.Decode(&getMediaResponse))

	return getMediaResponse.Items
}
