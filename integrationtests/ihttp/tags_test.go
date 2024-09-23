package ihttp

import (
	"encoding/json"
	"fmt"
	"media-nexus/adapters/primary/ahttp/ahmodel"
	"media-nexus/integrationtests"
	"media-nexus/model"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

type tagsE2ETestSuite struct {
	integrationtests.E2ETestSuite
}

func TestTags(t *testing.T) {
	suite.Run(t, &tagsE2ETestSuite{})
}

func (s *tagsE2ETestSuite) TestCreateTag() {
	ctx := s.Context()

	tagID := s.createTag(s.GenerateAlphanumeric(10))

	err := s.App().TagRepo().DeleteTags(ctx, []model.TagID{tagID})
	s.Require().NoError(err)
}

func (s *tagsE2ETestSuite) TestCreateTagIdempotency() {
	ctx := s.Context()

	tagName := s.GenerateAlphanumeric(10)
	tagID1 := s.createTag(tagName)

	defer func() {
		s.LogIfError(s.App().TagRepo().DeleteTags(ctx, []model.TagID{tagID1}), "delete tags")
	}()

	tagID2 := s.createTag(tagName)

	s.Equal(tagID1, tagID2)
}

func (s *tagsE2ETestSuite) TestListTags() {
	ctx := s.Context()

	var tagIds []model.TagID
	tagIds = append(tagIds, s.createTag(s.GenerateAlphanumeric(10)))
	tagIds = append(tagIds, s.createTag(s.GenerateAlphanumeric(10)))

	defer func() { s.LogIfError(s.App().TagRepo().DeleteTags(ctx, tagIds), "delete tags") }()

	storedTagIDs := s.listTags()
	for _, tagID := range tagIds {
		s.Contains(storedTagIDs, tagID)
	}
}

func (s *tagsE2ETestSuite) createTag(tagName string) model.TagID {
	reqStr := fmt.Sprintf(`{"name":"%v"}`, tagName)

	req, err := http.NewRequest(http.MethodPost, s.CreateServerURL("/tags"), strings.NewReader(reqStr))
	s.NoError(err)

	req.Header.Set("Content-Type", "application/json")

	response, err := s.Client().Do(req)
	s.NoError(err)

	defer response.Body.Close()

	s.Equal(http.StatusOK, response.StatusCode)

	var postTagsResponse ahmodel.PostTagsResponse
	decoder := json.NewDecoder(response.Body)
	s.NoError(decoder.Decode(&postTagsResponse))

	s.Require().NotEmpty(postTagsResponse.TagID)

	return model.TagID(postTagsResponse.TagID)
}

func (s *tagsE2ETestSuite) listTags() []model.TagID {
	req, err := http.NewRequest(http.MethodGet, s.CreateServerURL("/tags"), nil)
	s.NoError(err)

	req.Header.Set("Content-Type", "application/json")

	response, err := s.Client().Do(req)
	s.NoError(err)

	defer response.Body.Close()

	s.Equal(http.StatusOK, response.StatusCode)

	var postTagsResponse []*ahmodel.Tag
	decoder := json.NewDecoder(response.Body)
	s.NoError(decoder.Decode(&postTagsResponse))

	result := make([]model.TagID, 0, len(postTagsResponse))
	for _, t := range postTagsResponse {
		result = append(result, model.TagID(t.ID))
	}

	return result
}
