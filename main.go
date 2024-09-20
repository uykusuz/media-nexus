package main

import (
	"context"
	"flag"
	"os"

	"media-nexus/adapters/primary/ahttp"
	"media-nexus/adapters/secondary/aaws"
	"media-nexus/adapters/secondary/amongodb"
	"media-nexus/config"
	"media-nexus/errortypes"
	"media-nexus/logger"
	"media-nexus/services"
	"media-nexus/util"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var log = logger.NewLogger("main")

func main() {
	envPrefix := "MEDIANEXUS"

	err := run(log, envPrefix)
	if err != nil {
		log.Fatalf("%+v", err)
		os.Exit(1)
	}
}

func run(log logger.Logger, envPrefix string) error {
	ctx := util.WithLogger(context.Background(), log)

	flag.Parse()

	log.Infof("Reading config ...")

	if flag.NArg() < 1 {
		return errortypes.NewBadUserInput("Provide a path to the config file.")
	}

	config, err := config.ParseConfigurationFile(flag.Arg(0), envPrefix)
	if err != nil {
		return err
	}

	log.Info("setting up services ...")

	var runners []util.Runner

	awsConfig, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		return errortypes.NewIllegalStatef("failed to load aws config: %v", err)
	}

	s3Client := s3.NewFromConfig(awsConfig)
	presignClient := s3.NewPresignClient(s3Client)

	mongodbClient, err := amongodb.NewMongoDBClient(ctx, config.MongoDBURI)
	if err != nil {
		return errortypes.NewIllegalStatef("failed to create mongodb client: %v", err)
	}

	tags := amongodb.NewTagRepository(mongodbClient, config.MediaDatabase, config.MediaTagCollection)
	mediaMetadata, mediaMetadataRunner := amongodb.NewMediaMetadataRepository(
		mongodbClient,
		config.MediaDatabase,
		config.MediaMetadataCollection,
	)
	runners = append(runners, mediaMetadataRunner)

	media := aaws.NewMediaRepository(s3Client, presignClient, config.MediaBucket)
	mediaService := services.NewMediaService(
		tags,
		mediaMetadata,
		media,
		config.MediaUrlLifetime,
		config.IncompleteMediaMetadataLifetime,
	)

	log.Info("starting async tasks ...")
	for _, runner := range runners {
		go runner(ctx)
	}

	log.Info("starting API ...")
	return ahttp.StartAPI(log, config.HTTPPort, mediaService, tags)
}
