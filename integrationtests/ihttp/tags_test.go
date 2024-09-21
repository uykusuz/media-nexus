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

	tagName, err := randutil.Alphanumeric(10)
	s.Require().NoError(err)
	reqStr := fmt.Sprintf(`{"name":"%v"}`, tagName)

	req, err := http.NewRequest(http.MethodPost, s.CreateServerUrl("/tags"), strings.NewReader(reqStr))
	s.NoError(err)

	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	response, err := client.Do(req)
	s.NoError(err)

	defer response.Body.Close()

	s.Equal(http.StatusOK, response.StatusCode)

	var postTagsResponse ahmodel.PostTagsResponse
	decoder := json.NewDecoder(response.Body)
	s.NoError(decoder.Decode(&postTagsResponse))

	s.Require().NotEmpty(postTagsResponse.TagId)

	err = s.App().TagRepo().DeleteTag(ctx, model.TagId(postTagsResponse.TagId))
	s.Require().NoError(err)
}
