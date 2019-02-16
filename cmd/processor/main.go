package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"gitlab.com/Cacophony/Processor/pkg/processing"
	"gitlab.com/Cacophony/Processor/plugins"
	"gitlab.com/Cacophony/go-kit/api"
	"gitlab.com/Cacophony/go-kit/logging"
	"gitlab.com/Cacophony/go-kit/state"
	"go.uber.org/zap"
)

const (
	// ServiceName is the name of the service
	ServiceName = "processor"
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

	// init redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddress,
		Password: config.RedisPassword,
	})
	_, err = redisClient.Ping().Result()
	if err != nil {
		logger.Fatal("unable to connect to Redis",
			zap.Error(err),
		)
	}

	// init AMQP session
	amqpConnection, err := amqp.Dial(config.AMQPDSN)
	if err != nil {
		logger.Fatal("unable to initialise AMQP session",
			zap.Error(err),
		)
	}

	// init GORM
	gormDB, err := gorm.Open("postgres", config.DBDSN)
	if err != nil {
		logger.Fatal("unable to initialise GORM session",
			zap.Error(err),
		)
	}
	// gormDB.SetLogger(logger) TODO: write logger
	defer gormDB.Close()

	// init state
	stateClient := state.NewSate(redisClient)

	// init processor
	processor, err := processing.NewProcessor(
		logger.With(zap.String("feature", "processor")),
		ServiceName,
		amqpConnection,
		"cacophony",
		"cacophony.discord.#",
		gormDB,
		stateClient,
		config.ConcurrentProcessingLimit,
		config.ProcessingDeadline,
		config.DiscordTokens,
	)
	if err != nil {
		logger.Fatal("unable to initialise processor",
			zap.Error(err),
		)
	}

	// init http server
	httpRouter := api.NewRouter()
	httpServer := api.NewHTTPServer(config.Port, httpRouter)

	// init plugins
	plugins.StartPlugins(
		logger.With(zap.String("feature", "start_plugins")),
		gormDB,
	)

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
		zap.Int("concurrent_processing_limit", config.ConcurrentProcessingLimit),
		zap.String("environment", string(config.Environment)),
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

	// TODO: make sure processor is finished processing events before shutting down

	plugins.StopPlugins(
		logger.With(zap.String("feature", "stop_plugins")),
		gormDB,
	)

	err = httpServer.Shutdown(ctx)
	if err != nil {
		logger.Error("unable to shutdown HTTP Server",
			zap.Error(err),
		)
	}

}
