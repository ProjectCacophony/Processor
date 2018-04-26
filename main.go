package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/json-iterator/go"
	"gitlab.com/project-d-collab/SqsProcessor/modules"
	"gitlab.com/project-d-collab/dhelpers"
	dhelpersCache "gitlab.com/project-d-collab/dhelpers/cache"
	"gitlab.com/project-d-collab/dhelpers/components"
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

	// Setup all modules
	modules.Init()

	dhelpersCache.GetLogger().Infoln("Processor booting completed, took", time.Since(started).String())

	// bot run loop
	go func() {
		sqsClient := dhelpersCache.GetAwsSqsSession()
		logger := dhelpersCache.GetLogger()
		redisClient := dhelpersCache.GetRedisClient()

		for {
			result, err := sqsClient.ReceiveMessage(&sqs.ReceiveMessageInput{
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

	// Uninit all modules
	modules.Uninit()
}
