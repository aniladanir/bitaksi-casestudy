package main

import (
	"context"
	"flag"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aniladanir/bitaksi-casestudy/matching-api/internal/adapters/handlers/httphandler"
	"github.com/aniladanir/bitaksi-casestudy/matching-api/internal/adapters/locationfinder"
	"github.com/aniladanir/bitaksi-casestudy/matching-api/internal/core/services"
	"github.com/aniladanir/bitaksi-casestudy/matching-api/pkg/config"
	"github.com/aniladanir/bitaksi-casestudy/matching-api/pkg/log"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

var configFile = flag.String("config", "./config.yaml", "provide configuration file")

func main() {
	flag.Parse()

	httpHandler, err := InitializeComponents()
	if err != nil {
		log.Fatal("encountered error when initializing components", zap.Error(err))
	}

	ctx, cancel := context.WithCancel(context.Background())
	errGroup, errGroupCtx := errgroup.WithContext(ctx)

	// listen os signals and cancel the parent context if one is received.
	go ListenOsSignal(cancel)

	// start http handler
	errGroup.Go(func() error {
		err := httpHandler.Listen(config.GetHttpServerAddress())
		return fmt.Errorf("http handler failed listening: %v", err)
	})

	// gracefully shutdown application if context is canceled
	errGroup.Go(func() error {
		log.Info("matching-api is shutting down...")
		<-errGroupCtx.Done()

		// Shut down httplkerfd server
		if err := httpHandler.Shutdown(); err != nil {
			return err
		}

		return nil
	})

	// wait for graceful shutdown
	if err := errGroup.Wait(); err != nil {
		log.Error("encountered error on shutdown", zap.Error(err))
	}

	log.Info("greceful shutdown is complete")
}

// InitializeComponents initalizes all adapters needed by application
func InitializeComponents() (*httphandler.Handler, error) {
	config.Init(*configFile)

	// configure logrotate options
	logRotateCfg := log.RotateConfig{
		MaxSizeMB:   config.GetLogMaxSizeInMB(),
		MaxAgeDays:  config.GetLogMaxAgeInDays(),
		MaxBackups:  config.GetLogMaxBackups(),
		GzipArchive: config.GetLogGzipArchive(),
	}

	// create driver location api client
	driverLocationApiUrl, err := url.Parse(config.GetRemoteUrl("driverLocationApi"))
	if err != nil {
		return nil, fmt.Errorf("could not parse driver location api url: %w", err)
	}
	driverLocationApiClient := locationfinder.NewDriverLocationApiClient(
		*driverLocationApiUrl,
		config.GetRemoteVersion("driverLocationApi"),
		time.Duration(config.GetHttpClientTimeout())*time.Second,
	)

	// create http handler
	debug := config.IsDebug()
	httpHandler := httphandler.NewHandler(
		httphandler.ServerConfig{
			WriteTimeout: time.Duration(config.GetHttpWriteTimeout()) * time.Second,
			ReadTimeout:  time.Duration(config.GetHttpReadTimeout()) * time.Second,
			IdleTimeout:  time.Duration(config.GetHttpIdleTimeout()) * time.Second,
		},
		log.NewLoggerWithLogRotate(debug, config.GetLogFile(), logRotateCfg),
		log.NewLoggerWithLogRotate(debug, config.GetAccessLogFile(), logRotateCfg),
		services.NewDriverService(driverLocationApiClient),
		config.GetAPIVersion(),
	)

	return httpHandler, nil
}

func ListenOsSignal(onSignal func()) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan

	onSignal()
}
