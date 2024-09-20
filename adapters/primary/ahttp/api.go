package ahttp

import (
	"fmt"
	"media-nexus/logger"
	"media-nexus/ports"
	"media-nexus/services"
	"net/http"
)

func StartAPI(log logger.Logger, port int, mediaService services.MediaService, tags ports.TagRepository) error {
	mediaEndpoint := &mediaEndpoint{mediaService, log, 200}
	http.HandleFunc("/api/v1/media", mediaEndpoint.HandleMedia)

	tagsEndpoint := &tagsEndpoint{tags, log}
	http.HandleFunc("/api/v1/tags", tagsEndpoint.HandleTags)

	log.Infof("start serving http on port %v ...", port)
	return http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
}
