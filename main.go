package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"strings"

	"sync"

	"github.com/json-iterator/go"
	"gitlab.com/Cacophony/Processor/api"
	"gitlab.com/Cacophony/Processor/modules"
	"gitlab.com/Cacophony/dhelpers"
	"gitlab.com/Cacophony/dhelpers/cache"
	"gitlab.com/Cacophony/dhelpers/components"
)

var (
	started time.Time
)

func main() {
	started = time.Now()
	var err error

	// Set up components
	components.InitMetrics()
	components.InitLogger("Processor")
	err = components.InitSentry()
	dhelpers.CheckErr(err, "failed to initialise sentry")
	components.InitTranslator(nil)
	components.InitRedis()
	err = components.InitMongoDB()
	dhelpers.CheckErr(err, "failed to initialise mongodb")
	components.InitLastFm()
	err = components.InitTracer("Processor")
	dhelpers.CheckErr(err, "failed to initialise tracer")
	err = components.InitKafkaConsumerGroup()
	dhelpers.CheckErr(err, "failed to initialise kafka consumer group")

	// start api server
	apiServer := &http.Server{
		Addr: os.Getenv("API_ADDRESS"),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      api.New(),
	}
	go func() {
		apiServerListenAndServeErr := apiServer.ListenAndServe()
		if err != nil && !strings.Contains(err.Error(), "http: Server closed") {
			cache.GetLogger().Fatal(apiServerListenAndServeErr)
		}
	}()
	cache.GetLogger().Infoln("started API on", os.Getenv("API_ADDRESS"))

	// Setup all modules
	modules.Init()

	cache.GetLogger().Infoln("Processor booting completed, took", time.Since(started).String())

	// bot run loop
	go func() {
		consumer := cache.GetKafkaConsumerGroup()
		logger := cache.GetLogger()
		redisClient := cache.GetRedisClient()

		for event := range consumer.Messages() {
			//for event := range saramaPartitionConsumer.Messages() {
			// unpack the event data
			var eventContainer dhelpers.EventContainer
			err = jsoniter.Unmarshal(event.Value, &eventContainer)
			if err != nil {
				logger.Errorln("Message unmarshal error: ", err.Error())
				continue
			}
			// deduplication
			if !dhelpers.IsNewEvent(redisClient, "sqs-processor", eventContainer.Key) {
				continue
			}

			// send to modules
			modules.CallModules(eventContainer)
		}
	}()
	go func() {
		consumer := cache.GetKafkaConsumerGroup()
		logger := cache.GetLogger()

		for event := range consumer.Errors() {
			//for event := range saramaPartitionConsumer.Errors() {
			logger.WithError(event).Errorln("received error from Kafka Consumer for Partition")
		}
	}()

	// channel for bot shutdown
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// create sync.WaitGroup for all shutdown goroutines
	var exitGroup sync.WaitGroup

	// modules uninit goroutine
	exitGroup.Add(1)
	go func() {
		// uninit all modules
		cache.GetLogger().Infoln("Uniniting all modules…")
		modules.Uninit()
		cache.GetLogger().Infoln("Uninited all modules")
		exitGroup.Done()
	}()

	// API Server shutdown goroutine
	exitGroup.Add(1)
	go func() {
		// shutdown api server
		cache.GetLogger().Infoln("Shutting API server down…")
		err = apiServer.Shutdown(context.Background())
		dhelpers.LogError(err)
		cache.GetLogger().Infoln("Shut API server down")
		exitGroup.Done()
	}()

	// Kafka Consumer Group shutdown goroutine
	exitGroup.Add(1)
	go func() {
		// shutdown Kafka Consumer Group
		cache.GetLogger().Infoln("Shutting Kafka Consumer Group down…")
		err = cache.GetKafkaConsumerGroup().Close()
		dhelpers.LogError(err)
		cache.GetLogger().Infoln("Shut Kafka Consumer Group down")

		exitGroup.Done()
	}()

	// wait for all shutdown goroutines to finish and then close channel finished
	finished := make(chan bool)
	go func() {
		exitGroup.Wait()
		close(finished)
	}()

	// wait 60 second for everything to finish, or shut it down anyway
	select {
	case <-finished:
		cache.GetLogger().Infoln("shutdown successful")
	case <-time.After(60 * time.Second):
		cache.GetLogger().Infoln("forcing shutdown after 60 seconds")
	}
}
