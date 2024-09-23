package config

import (
	"time"

	"media-nexus/validation"
)

// Configuration defines options for the application.
type Configuration struct {
	BaseURL  string
	HTTPPort int

	MongoDBURI                      string
	MediaDatabase                   string
	MediaTagCollection              string
	MediaMetadataCollection         string
	MediaBucket                     string
	MediaBucketRegion               string
	GetMediaURLLifetime             time.Duration
	IncompleteMediaMetadataLifetime time.Duration
}

func NewConfiguration() Configuration {
	return Configuration{
		BaseURL:                         "http://localhost",
		HTTPPort:                        8081,
		MongoDBURI:                      "http://localhost:27017",
		MediaDatabase:                   "media",
		MediaTagCollection:              "tags",
		MediaMetadataCollection:         "media_metadata",
		MediaBucket:                     "hintergarten.de-media-nexus-media",
		MediaBucketRegion:               "eu-central-1",
		GetMediaURLLifetime:             15 * 60 * time.Second,
		IncompleteMediaMetadataLifetime: 60 * time.Second,
	}
}

func (c *Configuration) Validate() error {
	if err := validation.IsValidStringProperty("<root>", "baseUrl", c.BaseURL); err != nil {
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

	if err := validation.IsValidStringProperty("<root>", "mediaBucketRegion", c.MediaBucketRegion); err != nil {
		return err
	}

	return nil
}
