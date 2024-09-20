package httputils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"media-nexus/logger"

	"github.com/pkg/errors"
)

// RespondWithJSON uses the json encoder to write an interface to the http response with an optional status code.
func RespondWithJSON(
	statusCode int,
	i interface{},
	response http.ResponseWriter,
	logger logger.Logger,
	escapeHTML bool,
) {
	body, err := EncodeJSON(i, escapeHTML)
	if err != nil {
		RespondWithInternalError(response)
		logger.Error(err)
		return
	}

	response.Header().Set(HeaderContentType, ContentTypeJSON)
	response.WriteHeader(statusCode)

	_, err = response.Write(body)
	if err != nil {
		logger.Errorf("failed to write response body: %v", err)
	}
}

func EncodeJSON(i interface{}, escapeHTML bool) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.SetEscapeHTML(escapeHTML)

	err := encoder.Encode(i)
	if err != nil {
		return nil, errors.Wrap(err, "failed to encode json response")
	}

	return buffer.Bytes(), nil
}

func RespondWithError(
	response http.ResponseWriter,
	httpStatusCode int,
	format string,
	a ...interface{},
) {
	errorMessage := fmt.Sprintf(format, a...)
	http.Error(response, errorMessage, httpStatusCode)
}

func RespondWithInternalError(response http.ResponseWriter) int {
	RespondWithError(response, http.StatusInternalServerError, "internal server error")
	return http.StatusInternalServerError
}

func RespondWithBadParameter(response http.ResponseWriter, parameterName string, reason error) {
	RespondWithError(
		response,
		http.StatusBadRequest,
		"failed to parse '%v': %v",
		parameterName,
		reason,
	)
}
