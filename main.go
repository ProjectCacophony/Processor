package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis"
	"github.com/json-iterator/go"
	"gitlab.com/project-d-collab/dhelpers"
)

var (
	token           string
	awsRegion       string
	started         time.Time
	sqsClient       *sqs.SQS
	sqsQueueUrl     string
	dg              *discordgo.Session
	redisAddress    string
	redisClient     *redis.Client
	discordEndpoint string
)

func init() {
	// Parse command line flags (-t DISCORD_BOT_TOKEN -aws-region AWS_REGION -redis REDIS_ADDRESS -sqs SQS_QUEUE_URL -discord-endpoint DISCORD_ENDPOINT)
	flag.StringVar(&token, "t", "", "Discord Bot token")
	flag.StringVar(&awsRegion, "aws-region", "", "AWS Region")
	flag.StringVar(&redisAddress, "redis", "127.0.0.1:6379", "Redis Address")
	flag.StringVar(&sqsQueueUrl, "sqs", "", "SQS Queue Url")
	flag.StringVar(&discordEndpoint, "discord-endpoint", "https://discordapp.com/", "Discord Endpoint URL")
	flag.Parse()
	// overwrite with environment variables if set DISCORD_BOT_TOKEN=… AWS_REGION=… REDIS_ADDRESS=… SQS_QUEUE_URL=… DISCORD_ENDPOINT=…
	if os.Getenv("DISCORD_BOT_TOKEN") != "" {
		token = os.Getenv("DISCORD_BOT_TOKEN")
	}
	if os.Getenv("AWS_REGION") != "" {
		awsRegion = os.Getenv("AWS_REGION")
	}
	if os.Getenv("REDIS_ADDRESS") != "" {
		redisAddress = os.Getenv("REDIS_ADDRESS")
	}
	if os.Getenv("SQS_QUEUE_URL") != "" {
		sqsQueueUrl = os.Getenv("SQS_QUEUE_URL")
	}
	if os.Getenv("DISCORD_ENDPOINT") != "" {
		discordEndpoint = os.Getenv("DISCORD_ENDPOINT")
	}
}

func main() {
	started = time.Now()
	var err error
	// connect to aws
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	}))
	sqsClient = sqs.New(sess)

	// connect to redis
	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: "",
		DB:       0,
	})

	// create a new Discordgo Bot Client
	dhelpers.SetDiscordEndpoints(discordEndpoint)
	fmt.Println("Set Discord Endpoint API URL to", discordgo.EndpointAPI)
	fmt.Println("Connecting to Discord, token Length:", len(token))
	dg, err = discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err.Error())
		return
	}

	for {
		result, err := sqsClient.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl:              aws.String(sqsQueueUrl),
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
				fmt.Println(err.Error())
				continue
			}
			// deduplication
			if !dhelpers.IsNewEvent(redisClient, "sqs-processor", eventContainer.Key) {
				return
			}

			receivedAt := time.Now()

			fmt.Println(eventContainer.Key)
			switch eventContainer.Type {
			case dhelpers.MessageCreateEventType:
				if strings.Contains(eventContainer.MessageCreate.Content, "ping") {
					_, err = dg.ChannelMessageSend(eventContainer.MessageCreate.ChannelID, "pong, Gateway => SqsProcessor: "+receivedAt.Sub(eventContainer.ReceivedAt).String())
					if err != nil {
						fmt.Println(err.Error())
					}
				}
			}
		}
	}
}
