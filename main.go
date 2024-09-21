package main

import (
	"os"

	"media-nexus/app"
	"media-nexus/config"
	"media-nexus/logger"
)

func main() {
	var log = logger.NewLogger("main")

	if err := run(log); err != nil {
		log.Fatalf("%+v", err)
		os.Exit(1)
	}
}

func run(log logger.Logger) error {
	log.Infof("Reading config ...")

	config, err := config.LoadConfiguration()
	if err != nil {
		return err
	}

	application := app.NewApp(log, config)

	if err := application.Setup(); err != nil {
		return err
	}

	return application.Run()
}
