package httputils

import (
	"encoding/json"
	"io"
	"net/url"
	"strconv"

	"media-nexus/errortypes"
)

func ParseInt32QueryParameter(queryValues url.Values, paramName string) (int32, error) {
	param := queryValues.Get(paramName)

	if len(param) < 1 {
		return 0, errortypes.NewBadUserInputf("param %v empty", paramName)
	}

	i, err := strconv.ParseInt(param, 10, 32)
	if err != nil {
		return 0, errortypes.NewBadUserInputf("failed to parse %v as integer: %v", paramName, err)
	}

	return int32(i), nil
}

func ParseJSONRequestBody(body io.Reader, output interface{}) error {
	if err := json.NewDecoder(body).Decode(output); err != nil {
		return errortypes.NewBadUserInputf("failed to parse request body as JSON: %v", err)
	}

	return nil
}
