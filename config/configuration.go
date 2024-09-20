package config

import (
	"time"

	"media-nexus/validation"
)

// Configuration defines options for the application.
type Configuration struct {
	HTTPPort int

	MongoDBURI                      string
	MediaDatabase                   string
	MediaTagCollection              string
	MediaMetadataCollection         string
	MediaBucket                     string
	MediaUrlLifetime                time.Duration
	IncompleteMediaMetadataLifetime time.Duration
}

func NewConfiguration() Configuration {
	return Configuration{
		HTTPPort:                8081,
		MongoDBURI:              "http://localhost:27017",
		MediaDatabase:           "media",
		MediaTagCollection:      "tags",
		MediaMetadataCollection: "media_metadata",
		MediaBucket:             "hintergarten.de-media-nexus-media",
	}
}

func (c *Configuration) Validate() error {
	if err := validation.IsValidPortProperty("<root>", "httpPort", c.HTTPPort); err != nil {
		return err
	}

	return nil
}
