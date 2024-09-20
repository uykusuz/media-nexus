package httputils

import (
	"net/http"

	"media-nexus/errortypes"
	"media-nexus/logger"

	"github.com/pkg/errors"
)

type LogSomething int

const (
	LogNothing LogSomething = iota
	LogInfo
	LogError
)

func HandleError(err error, response http.ResponseWriter, logger logger.Logger) bool {
	if err == nil {
		return false
	}

	var code int
	log := LogError
	respondWithError := true
	responseMessage := err.Error()

	switch errors.Cause(err).(type) {
	case errortypes.BadUserInput:
		code = http.StatusBadRequest
		log = LogInfo
	case errortypes.InputOutputError:
		code = http.StatusInternalServerError
	case errortypes.ResourceAlreadyExists:
		code = http.StatusConflict
	default:
		code = RespondWithInternalError(response)
		respondWithError = false
	}

	switch log {
	case LogInfo:
		logger.Infof("%v\n", err)
	case LogError:
		logger.Errorf("%v\n", err)
	}

	if respondWithError {
		RespondWithError(response, code, responseMessage)
	}

	return true
}
