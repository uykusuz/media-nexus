package config

import (
	"time"

	"media-nexus/validation"
)

// Configuration defines options for the application.
type Configuration struct {
	BaseUrl string
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
		BaseUrl: "http://localhost",
		HTTPPort:                8081,
		MongoDBURI:              "http://localhost:27017",
		MediaDatabase:           "media",
		MediaTagCollection:      "tags",
		MediaMetadataCollection: "media_metadata",
		MediaBucket:             "hintergarten.de-media-nexus-media",
	}
}

func (c *Configuration) Validate() error {
	if err := validation.IsValidStringProperty("<root>", "baseUrl", c.BaseUrl); err != nil {
		return err
	}

	if err := validation.IsValidPortProperty("<root>", "httpPort", c.HTTPPort); err != nil {
		return err
	}

	if err := validation.IsValidStringProperty("<root>", "mongDbUri", c.MongoDBURI); err != nil {
		return err
	}

	if err := validation.IsValidStringProperty("<root>", "mediaDatabase", c.MediaDatabase); err != nil {
		return err
	}

	if err := validation.IsValidStringProperty("<root>", "mediaTagCollection", c.MediaTagCollection); err != nil {
		return err
	}

	if err := validation.IsValidStringProperty("<root>", "mediaMetadataCollection", c.MediaMetadataCollection); err != nil {
		return err
	}

	if err := validation.IsValidStringProperty("<root>", "mediaBucket", c.MediaBucket); err != nil {
		return err
	}

	return nil
}
