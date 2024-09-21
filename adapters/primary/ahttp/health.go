package ahttp

import (
	"media-nexus/logger"
	"net/http"
)

type healthEndpoint struct {
	log logger.Logger
}

func newHealthEndpoint(
	log logger.Logger,
) (*healthEndpoint, error) {
	return &healthEndpoint{log}, nil
}

// GetHealthLive godoc
//
//	@Summary		Retrieve live health status of server
//	@Tags			tags
//	@Produce		json
//	@Success		204
//	@Router			/api/v1/health/live [get]
func (e *healthEndpoint) GetHealthLive(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

// GetHealthReady godoc
//
//	@Summary		Retrieve ready health status of server
//	@Tags			tags
//	@Produce		json
//	@Success		204
//	@Router			/api/v1/health/ready [get]
func (e *healthEndpoint) GetHealthReady(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
