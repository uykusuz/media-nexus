package ahttp

import (
	"context"
	"fmt"
	"media-nexus/logger"
	"media-nexus/ports"
	"media-nexus/services"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "media-nexus/docs"

	"github.com/gorilla/mux"
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
	baseURL string,
	port int,
	mediaService services.MediaService,
	tags ports.TagRepository,
) error {
	r := mux.NewRouter()

	healthEndpoint := &healthEndpoint{log}
	r.HandleFunc("/api/v1/health/live", healthEndpoint.GetHealthLive).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/health/ready", healthEndpoint.GetHealthReady).Methods(http.MethodGet)

	mediaEndpoint := &mediaEndpoint{mediaService, log, 200, 500, 200}
	r.HandleFunc("/api/v1/media", mediaEndpoint.GetMedia).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/media", mediaEndpoint.CreateMedia).Methods(http.MethodPost)

	tagsEndpoint := &tagsEndpoint{tags, log, 500}
	r.HandleFunc("/api/v1/tags", tagsEndpoint.ListTags).Methods(http.MethodGet)
	r.HandleFunc("/api/v1/tags", tagsEndpoint.CreateTag).Methods(http.MethodPost)

	r.PathPrefix("/swagger").Handler(createSwaggerHandler(baseURL, port)).Methods(http.MethodGet)

	srv := &http.Server{
		Addr:              fmt.Sprintf(":%v", port),
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Infof("start serving http on port %v ...", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Errorf("ListenAndServe error: %v\n", err)
		}
	}()

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Infof("Shutting down server...")
	return srv.Shutdown(ctx)
}

func createSwaggerHandler(baseURL string, port int) http.HandlerFunc {
	url := fmt.Sprintf("%v:%v/swagger/doc.json", baseURL, port)
	return httpSwagger.Handler(httpSwagger.URL(url))
}
