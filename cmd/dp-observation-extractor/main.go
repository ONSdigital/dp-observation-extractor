package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/ONSdigital/dp-observation-extractor/config"
	"github.com/ONSdigital/dp-observation-extractor/initialise"
	"github.com/ONSdigital/dp-observation-extractor/service"
	"github.com/ONSdigital/log.go/log"
)

var (
	// BuildTime represents the time in which the service was built
	BuildTime string
	// GitCommit represents the commit (SHA-1) hash of the service that is running
	GitCommit string
	// Version represents the version of the service that is running
	Version string
)

func main() {

	log.Namespace = "dp-observation-extractor"
	ctx := context.Background()

	config, err := config.Get()
	if err != nil {
		log.Event(ctx, "error getting config", log.ERROR, log.Error(err))
		os.Exit(1)
	}

	// sensitive fields are omitted from config.String().
	log.Event(ctx, "config on startup", log.INFO, log.Data{"config": config})

	// a channel used to signal a graceful exit is required.
	errorChannel := make(chan error)

	// Signal channel to know if SIGTERM is triggered
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	// serviceList keeps track of what dependency services have been initialised
	serviceList := initialise.ExternalServiceList{}

	err = service.Run(ctx, config, serviceList, signals, errorChannel, BuildTime, GitCommit, Version)
	if err != nil {
		log.Event(ctx, "error running service", log.ERROR, log.Error(err))
		os.Exit(1)
	}
}
