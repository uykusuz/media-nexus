package ahttp

import (
	"fmt"
	"media-nexus/logger"
	"media-nexus/ports"
	"media-nexus/services"
	"net/http"

	_ "media-nexus/docs"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

//	@title			media-nexus API
//	@version		1.0
//	@description	API to serve media assets

//	@contact.name	Boris Brönner
//	@contact.url	https://hintergarten.de
//	@contact.email	broenner@hintergarten.de

//	@license.name	MIT

// @host		localhost:8081
// @BasePath	/api/v1
func StartAPI(
	log logger.Logger,
	baseUrl string,
	port int,
	mediaService services.MediaService,
	tags ports.TagRepository,
) error {
	mediaEndpoint := &mediaEndpoint{mediaService, log, 200}
	http.HandleFunc("/api/v1/media", mediaEndpoint.HandleMedia)

	tagsEndpoint := &tagsEndpoint{tags, log}
	http.HandleFunc("/api/v1/tags", tagsEndpoint.HandleTags)

	http.HandleFunc("/swagger/", func(w http.ResponseWriter, r *http.Request) {
		swaggerHandler(baseUrl, port, w, r)
	})

	log.Infof("start serving http on port %v ...", port)
	return http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
}

func swaggerHandler(baseUrl string, port int, w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%v:%v/swagger/doc.json", baseUrl, port)
	httpSwagger.Handler(httpSwagger.URL(url))(w, r)
}
