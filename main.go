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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/json-iterator/go"
	"gitlab.com/Cacophony/SqsProcessor/api"
	"gitlab.com/Cacophony/SqsProcessor/modules"
	"gitlab.com/Cacophony/dhelpers"
	"gitlab.com/Cacophony/dhelpers/cache"
	"gitlab.com/Cacophony/dhelpers/components"
)

var (
	started     time.Time
	sqsQueueURL string
)

func init() {
	// parse environment variables
	sqsQueueURL = os.Getenv("SQS_QUEUE_URL")
}

func main() {
	started = time.Now()
	var err error

	// Set up components
	components.InitMetrics()
	components.InitLogger("SqsProcessor")
	err = components.InitSentry()
	dhelpers.CheckErr(err)
	components.InitTranslator(nil)
	components.InitRedis()
	err = components.InitMongoDB()
	dhelpers.CheckErr(err)
	err = components.InitAwsSqs()
	dhelpers.CheckErr(err)
	components.InitLastFm()
	err = components.InitTracer("SqsProcessor")
	dhelpers.CheckErr(err)

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

	cache.GetLogger().Infoln("SqsProcessor booting completed, took", time.Since(started).String())

	// bot run loop
	go func() {
		sqsClient := cache.GetAwsSqsSession()
		logger := cache.GetLogger()
		redisClient := cache.GetRedisClient()
		var result *sqs.ReceiveMessageOutput

		for {
			result, err = sqsClient.ReceiveMessage(&sqs.ReceiveMessageInput{
				QueueUrl:              aws.String(sqsQueueURL),
				MaxNumberOfMessages:   aws.Int64(10),
				MessageAttributeNames: aws.StringSlice([]string{}),
				WaitTimeSeconds:       aws.Int64(20),
				VisibilityTimeout:     aws.Int64(60 * 60 * 12),
			})
			if err != nil {
				panic(err)
			}

			for _, message := range result.Messages {
				// unpack the event data
				var eventContainer dhelpers.EventContainer
				err = jsoniter.Unmarshal([]byte(*message.Body), &eventContainer)
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
