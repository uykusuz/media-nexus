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
	"go.step.sm/crypto/randutil"
)

type tagsE2ETestSuite struct {
	integrationtests.E2ETestSuite
}

func TestTags(t *testing.T) {
	suite.Run(t, &tagsE2ETestSuite{})
}

func (s *tagsE2ETestSuite) TestCreateTag() {
	ctx := s.Context()

	tagId := s.createTag(s.generateTagname())

	err := s.App().TagRepo().DeleteTags(ctx, []model.TagId{tagId})
	s.Require().NoError(err)
}

func (s *tagsE2ETestSuite) TestCreateTagIdempotency() {
	ctx := s.Context()

	tagName := s.generateTagname()
	tagId1 := s.createTag(tagName)

	defer func() {
		s.App().TagRepo().DeleteTags(ctx, []model.TagId{tagId1})
	}()

	tagId2 := s.createTag(tagName)

	s.Equal(tagId1, tagId2)
}

func (s *tagsE2ETestSuite) TestListTags() {
	ctx := s.Context()

	var tagIds []model.TagId
	tagIds = append(tagIds, s.createTag(s.generateTagname()))
	tagIds = append(tagIds, s.createTag(s.generateTagname()))

	defer func() {
		s.App().TagRepo().DeleteTags(ctx, tagIds)
	}()

	storedTagIds := s.listTags()
	for _, tagId := range tagIds {
		s.Contains(storedTagIds, tagId)
	}
}

func (s *tagsE2ETestSuite) generateTagname() string {
	tagName, err := randutil.Alphanumeric(10)
	s.Require().NoError(err)
	return tagName
}

func (s *tagsE2ETestSuite) createTag(tagName string) model.TagId {
	reqStr := fmt.Sprintf(`{"name":"%v"}`, tagName)

	req, err := http.NewRequest(http.MethodPost, s.CreateServerUrl("/tags"), strings.NewReader(reqStr))
	s.NoError(err)

	req.Header.Set("Content-Type", "application/json")

	response, err := s.Client().Do(req)
	s.NoError(err)

	defer response.Body.Close()

	s.Equal(http.StatusOK, response.StatusCode)

	var postTagsResponse ahmodel.PostTagsResponse
	decoder := json.NewDecoder(response.Body)
	s.NoError(decoder.Decode(&postTagsResponse))

	s.Require().NotEmpty(postTagsResponse.TagId)

	return model.TagId(postTagsResponse.TagId)
}

func (s *tagsE2ETestSuite) listTags() []model.TagId {
	req, err := http.NewRequest(http.MethodGet, s.CreateServerUrl("/tags"), nil)
	s.NoError(err)

	req.Header.Set("Content-Type", "application/json")

	response, err := s.Client().Do(req)
	s.NoError(err)

	defer response.Body.Close()

	s.Equal(http.StatusOK, response.StatusCode)

	var postTagsResponse []*ahmodel.Tag
	decoder := json.NewDecoder(response.Body)
	s.NoError(decoder.Decode(&postTagsResponse))

	result := make([]model.TagId, 0, len(postTagsResponse))
	for _, t := range postTagsResponse {
		result = append(result, model.TagId(t.Id))
	}

	return result
}
