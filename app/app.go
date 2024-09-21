package app

import (
	"context"

	"media-nexus/adapters/primary/ahttp"
	"media-nexus/adapters/secondary/aaws"
	"media-nexus/adapters/secondary/amongodb"
	"media-nexus/config"
	"media-nexus/errortypes"
	"media-nexus/logger"
	"media-nexus/ports"
	"media-nexus/services"
	"media-nexus/util"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type App interface {
	Setup() error
	Run() error

	TagRepo() ports.TagRepository
}

func NewApp(log logger.Logger, config *config.Configuration) App {
	return &app{
		log:    log,
		config: config,
	}
}

type app struct {
	log    logger.Logger
	config *config.Configuration

	runners      []util.Runner
	mediaService services.MediaService
	tagRepo      ports.TagRepository
}

func (a *app) Setup() error {
	ctx := util.WithLogger(context.Background(), a.log)

	a.log.Info("setting up services ...")

	awsConfig, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		return errortypes.NewIllegalStatef("failed to load aws config: %v", err)
	}

	s3Client := s3.NewFromConfig(awsConfig, aaws.WithRegion(a.config.MediaBucketRegion))
	presignClient := s3.NewPresignClient(s3Client)

	mongodbClient, err := amongodb.NewMongoDBClient(ctx, a.config.MongoDBURI)
	if err != nil {
		return errortypes.NewIllegalStatef("failed to create mongodb client: %v", err)
	}

	a.tagRepo = amongodb.NewTagRepository(mongodbClient, a.config.MediaDatabase, a.config.MediaTagCollection)
	mediaMetadata, mediaMetadataRunner := amongodb.NewMediaMetadataRepository(
		mongodbClient,
		a.config.MediaDatabase,
		a.config.MediaMetadataCollection,
		a.config.IncompleteMediaMetadataLifetime,
	)
	a.runners = append(a.runners, mediaMetadataRunner)

	media := aaws.NewMediaRepository(s3Client, presignClient, a.config.MediaBucket)
	a.mediaService = services.NewMediaService(
		a.tagRepo,
		mediaMetadata,
		media,
		a.config.GetMediaUrlLifetime,
		a.config.IncompleteMediaMetadataLifetime,
	)

	return nil
}

func (a *app) Run() error {
	ctx := util.WithLogger(context.Background(), a.log)

	a.log.Info("starting async tasks ...")
	for _, runner := range a.runners {
		go runner(ctx)
	}

	a.log.Info("starting API ...")
	return ahttp.StartAPI(a.log, a.config.BaseUrl, a.config.HTTPPort, a.mediaService, a.tagRepo)
}

func (a *app) TagRepo() ports.TagRepository {
	return a.tagRepo
}
