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

	tagIds := s.createTags(ctx, 2)
	defer s.App().TagRepo().DeleteTags(ctx, tagIds)

	mediaId := s.createMedia(s.GenerateAlphanumeric(10), tagIds, "./../assets/test.png")

	s.App().MediaMetadataRepo().DeleteAll(ctx, []model.MediaId{mediaId})
	s.App().MediaRepo().DeleteAll(ctx, []string{mediaId})
}

func (s *mediaE2ETestSuite) TestCreateMediaIdempotency() {
	ctx := s.Context()

	tagIds := s.createTags(ctx, 2)
	defer s.App().TagRepo().DeleteTags(ctx, tagIds)

	name := s.GenerateAlphanumeric(10)
	mediaId := s.createMedia(name, tagIds, "./../assets/test.png")
	defer s.App().MediaMetadataRepo().DeleteAll(ctx, []model.MediaId{mediaId})
	defer s.App().MediaRepo().DeleteAll(ctx, []string{mediaId})

	mediaId2 := s.createMedia(name, tagIds, "./../assets/test.png")

	if !s.Equal(mediaId, mediaId2) {
		s.App().MediaMetadataRepo().DeleteAll(ctx, []model.MediaId{mediaId2})
		s.App().MediaRepo().DeleteAll(ctx, []string{mediaId2})
	}
}

func (s *mediaE2ETestSuite) TestGetMediaByTagIds() {
	ctx := s.Context()

	tagIds := s.createTags(ctx, 2)
	defer s.App().TagRepo().DeleteTags(ctx, tagIds)

	mediaId := s.createMedia(s.GenerateAlphanumeric(10), tagIds, "./../assets/test.png")
	mediaId2 := s.createMedia(s.GenerateAlphanumeric(10), []model.TagId{tagIds[0]}, "./../assets/test2.png")

	mediaIds := []model.MediaId{mediaId, mediaId2}

	defer s.App().MediaMetadataRepo().DeleteAll(ctx, mediaIds)
	defer s.App().MediaRepo().DeleteAll(ctx, mediaIds)

	mediaItems := s.getMedia(tagIds[0])
	s.Equal(2, len(mediaItems))

	mediaItems = s.getMedia(tagIds[1])
	s.Equal(1, len(mediaItems))
}

func (s *mediaE2ETestSuite) createTags(ctx context.Context, count int) []model.TagId {
	var tagIds []model.TagId

	for i := 0; i < count; i++ {
		tagId, err := s.App().TagRepo().CreateTag(ctx, s.GenerateAlphanumeric(10))
		s.Require().NoError(err)
		tagIds = append(tagIds, tagId)
	}

	return tagIds
}

func (s *mediaE2ETestSuite) createMedia(name string, tagIds []model.TagId, filePath string) model.MediaId {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	err := writer.WriteField("name", name)
	s.Require().NoError(err)

	for _, tagId := range tagIds {
		err := writer.WriteField("tag_ids[]", tagId)
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

	req, err := http.NewRequest("POST", s.CreateServerUrl("/media"), &body)
	s.Require().NoError(err)

	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := s.Client().Do(req)
	s.Require().NoError(err)
	defer resp.Body.Close()

	s.Equal(http.StatusOK, resp.StatusCode)

	var postMediaResponse ahmodel.PostMediaResponse
	decoder := json.NewDecoder(resp.Body)
	s.NoError(decoder.Decode(&postMediaResponse))

	s.Require().NotEmpty(postMediaResponse.MediaId)

	return postMediaResponse.MediaId
}

func (s *mediaE2ETestSuite) getMedia(tagId model.TagId) []*ahmodel.MediaItem {
	url := s.CreateServerUrl("/media?tag_id=%v", tagId)
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
