package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"strings"

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

	// shutdown api server
	apiServerShutdownContext, apiServerCancel := context.WithTimeout(context.Background(), time.Second*15)
	defer apiServerCancel()
	err = apiServer.Shutdown(apiServerShutdownContext)
	dhelpers.LogError(err)

	// Uninit all modules
	modules.Uninit()
}
