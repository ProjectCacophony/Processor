package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	cacophonyConfig "gitlab.com/Cacophony/go-kit/config"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/paginator"

	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"gitlab.com/Cacophony/Processor/plugins"
	"gitlab.com/Cacophony/go-kit/api"
	"gitlab.com/Cacophony/go-kit/errortracking"
	"gitlab.com/Cacophony/go-kit/featureflag"
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
	config.FeatureFlag.Environment = config.ClusterEnvironment
	config.ErrorTracking.Version = config.Hash
	config.ErrorTracking.Environment = config.ClusterEnvironment

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

	// init raven
	err = errortracking.Init(&config.ErrorTracking)
	if err != nil {
		logger.Error("unable to initialise errortracking",
			zap.Error(err),
		)
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

	// init GORM
	gormDB, err := gorm.Open("postgres", config.DBDSN)
	if err != nil {
		logger.Fatal("unable to initialise GORM session",
			zap.Error(err),
		)
	}
	// gormDB.SetLogger(logger) TODO: write logger
	defer gormDB.Close()

	// init cacophony config
	err = cacophonyConfig.InitConfig(gormDB)
	if err != nil {
		logger.Fatal("unable to initialise Cacophony Config",
			zap.Error(err),
		)
	}

	// init state
	botIDs := make([]string, len(config.DiscordTokens))
	var i int
	for botID := range config.DiscordTokens {
		botIDs[i] = botID
		i++
	}
	stateClient := state.NewSate(redisClient, botIDs)

	// init feature flagger
	featureFlagger, err := featureflag.New(&config.FeatureFlag)
	if err != nil {
		logger.Fatal("unable to initialise feature flagger",
			zap.Error(err),
		)
	}

	// init paginator
	paginatorClient, err := paginator.NewPaginator(
		logger.With(zap.String("feature", "paginator")),
		redisClient,
		stateClient,
		config.DiscordTokens,
	)
	if err != nil {
		logger.Fatal("unable to initialise paginator",
			zap.Error(err),
		)
	}

	// create handler
	handler := handle(
		logger,
		gormDB,
		stateClient,
		config.DiscordTokens,
		featureFlagger,
		redisClient,
		paginatorClient,
		config.ProcessingDeadline,
	)

	// init processor
	processor, err := events.NewProcessor(
		logger.With(zap.String("feature", "processor")),
		ServiceName,
		config.AMQPDSN,
		"cacophony",
		"cacophony.discord.#",
		config.ConcurrentProcessingLimit,
		handler,
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
		redisClient,
		config.DiscordTokens,
		stateClient,
		featureFlagger,
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

	// TODO: make sure processor is finished processing events before shutting down

	plugins.StopPlugins(
		logger.With(zap.String("feature", "stop_plugins")),
		gormDB,
		redisClient,
		config.DiscordTokens,
		stateClient,
		featureFlagger,
	)

	err = httpServer.Shutdown(ctx)
	if err != nil {
		logger.Error("unable to shutdown HTTP Server",
			zap.Error(err),
		)
	}

}
