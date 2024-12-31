package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aniladanir/bitaksi-casestudy/driver-location-api/internal/adapters/handlers/httphandler"
	"github.com/aniladanir/bitaksi-casestudy/driver-location-api/internal/adapters/importer"
	"github.com/aniladanir/bitaksi-casestudy/driver-location-api/internal/adapters/repositories"
	"github.com/aniladanir/bitaksi-casestudy/driver-location-api/internal/adapters/repositories/mongodb"
	"github.com/aniladanir/bitaksi-casestudy/driver-location-api/internal/core/services"
	"github.com/aniladanir/bitaksi-casestudy/shared/config"
	"github.com/aniladanir/bitaksi-casestudy/shared/log"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

var configFile = flag.String("config", "./config.yaml", "provide configuration file")
var coordinatesFile = flag.String("coordinates", "./coordinates.csv", "provide coordinates csv file")

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	errGroup, errGroupCtx := errgroup.WithContext(ctx)

	// listen os signals and cancel the parent context if one is received.
	go ListenOsSignal(cancel)

	httpHandler, err := InitializeComponents(errGroupCtx)
	if err != nil {
		log.Fatal("encountered error when initializing components", zap.Error(err))
	}

	// start http handler
	errGroup.Go(func() error {
		err := httpHandler.Listen(config.GetHttpServerAddress())
		return fmt.Errorf("http handler failed listening: %v", err)
	})

	// gracefully shutdown application if context is canceled
	errGroup.Go(func() error {
		<-errGroupCtx.Done()
		log.Info("driver-location-api is shutting down...")

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
func InitializeComponents(ctx context.Context) (*httphandler.Handler, error) {
	config.Init(*configFile)

	// configure logrotate options
	logRotateCfg := log.RotateConfig{
		MaxSizeMB:   config.GetLogMaxSizeInMB(),
		MaxAgeDays:  config.GetLogMaxAgeInDays(),
		MaxBackups:  config.GetLogMaxBackups(),
		GzipArchive: config.GetLogGzipArchive(),
	}

	// initialize mongo db client
	mongoClient, err := mongodb.NewMongoClient(config.GetDBConnectionString())
	if err != nil {
		return nil, err
	}
	mongoDB, err := mongodb.CreateDatabase(ctx, mongoClient, config.GetDBName())
	if err != nil {
		return nil, err
	}

	// create repositories
	locationRepo := repositories.NewLocationRepository(mongoDB)

	// create importer
	csvImporter := importer.NewCsvImporter(locationRepo)

	// create services
	locationService := services.NewLocationService(locationRepo, csvImporter)

	// create loggers
	debug := config.IsDebug()
	appLogger := log.NewLoggerWithLogRotate(debug, config.GetLogFile(), logRotateCfg)
	acccessLogger := log.NewLoggerWithLogRotate(debug, config.GetAccessLogFile(), logRotateCfg)

	// create handlers
	httpHandler := httphandler.NewHandler(
		httphandler.ServerConfig{
			WriteTimeout: time.Duration(config.GetHttpWriteTimeout()) * time.Second,
			ReadTimeout:  time.Duration(config.GetHttpReadTimeout()) * time.Second,
			IdleTimeout:  time.Duration(config.GetHttpIdleTimeout()) * time.Second,
		},
		appLogger,
		acccessLogger,
		locationService,
		config.GetAPIVersion(),
	)

	// import initial coordinates
	file, err := os.Open(*coordinatesFile)
	if err != nil {
		appLogger.Error("could not open file", zap.String("file", *coordinatesFile), zap.Error(err))
		return nil, err
	}
	defer file.Close()

	if err := locationService.ImportLocation(ctx, file); err != nil {
		appLogger.Error("could not import coordinates", zap.Error(err))
		return nil, err
	}

	return httpHandler, nil
}

func ListenOsSignal(onSignal func()) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan

	onSignal()
}
