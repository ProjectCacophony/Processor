package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gitlab.com/Cacophony/go-kit/api"

	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"gitlab.com/Cacophony/Processor/pkg/processing"
	"gitlab.com/Cacophony/go-kit/logging"
	"go.uber.org/zap"
)

const (
	// ServiceName is the name of the service
	ServiceName = "gateway"
)

func main() {
	// init config
	var config config
	err := envconfig.Process("", &config)
	if err != nil {
		panic(errors.Wrap(err, "unable to load configuration"))
	}

	// init logger
	logger, err := logging.NewLogger(
		config.Environment,
		ServiceName,
		config.LoggingDiscordWebhook,
		&http.Client{
			Timeout: 10 * time.Second,
		},
	)
	if err != nil {
		panic(errors.Wrap(err, "unable to initialise launcher"))
	}

	// init AMQP session
	amqpConnection, err := amqp.Dial(config.AMQPDSN)
	if err != nil {
		logger.Fatal("unable to initialise AMQP session",
			zap.Error(err),
		)
	}

	// init processor
	processor, err := processing.NewProcessor(
		logger.With(zap.String("feature", "processor")),
		ServiceName,
		amqpConnection,
		"cacophony",
		"cacophony.discord.#",
	)
	if err != nil {
		logger.Fatal("unable to initialise processor",
			zap.Error(err),
		)
	}

	// init http server
	httpRouter := api.NewRouter()
	httpServer := api.NewHTTPServer(config.Port, httpRouter)

	// start processor
	go func() {
		err := processor.Start()
		if err != nil {
			logger.Fatal("processor error received", zap.Error(err))
		}
	}()

	// start http server
	go func() {
		err := httpServer.ListenAndServe()
		if err != http.ErrServerClosed {
			logger.Fatal("http server error",
				zap.Error(err),
				zap.String("feature", "http-server"),
			)
		}
	}()

	logger.Info("service is running",
		zap.Int("port", config.Port),
	)

	// wait for CTRL+C to stop the service
	quitChannel := make(chan os.Signal, 1)
	signal.Notify(quitChannel, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-quitChannel

	// shutdown features

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	err = amqpConnection.Close()
	if err != nil {
		logger.Error("unable to close AMQP session",
			zap.Error(err),
		)
	}

	err = httpServer.Shutdown(ctx)
	if err != nil {
		logger.Error("unable to shutdown HTTP Server",
			zap.Error(err),
		)
	}

}
