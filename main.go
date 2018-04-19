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
	"github.com/json-iterator/go"
	"gitlab.com/project-d-collab/dhelpers"
)

var (
	token       string
	awsRegion   string
	started     time.Time
	sqsClient   *sqs.SQS
	sqsQueueUrl string
	dg          *discordgo.Session
)

func init() {
	// Parse command line flags (-t DISCORD_BOT_TOKEN -aws-region AWS_REGION -sqs SQS_QUEUE_URL)
	flag.StringVar(&token, "t", "", "Discord Bot token")
	flag.StringVar(&awsRegion, "aws-region", "", "AWS Region")
	flag.StringVar(&sqsQueueUrl, "sqs", "", "SQS Queue Url")
	flag.Parse()
	// overwrite with environment variables if set DISCORD_BOT_TOKEN=… AWS_REGION=… REDIS_ADDRESS=… SQS_QUEUE_URL=…
	if os.Getenv("DISCORD_BOT_TOKEN") != "" {
		token = os.Getenv("DISCORD_BOT_TOKEN")
	}
	if os.Getenv("AWS_REGION") != "" {
		awsRegion = os.Getenv("AWS_REGION")
	}
	if os.Getenv("SQS_QUEUE_URL") != "" {
		sqsQueueUrl = os.Getenv("SQS_QUEUE_URL")
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

	// create a new Discordgo Bot Client
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
			MessageAttributeNames: aws.StringSlice([]string{"All"}),
			WaitTimeSeconds:       aws.Int64(20),
			VisibilityTimeout:     aws.Int64(60 * 60 * 12),
		})
		if err != nil {
			panic(err)
		}

		for _, message := range result.Messages {
			// pack the event data
			var eventContainer dhelpers.EventContainer
			err = jsoniter.Unmarshal([]byte(*message.Body), &eventContainer)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			fmt.Println(eventContainer.Key)
			switch eventContainer.Type {
			case dhelpers.MessageCreateEventType:
				if strings.Contains(eventContainer.MessageCreate.Content, "ping") {
					dg.ChannelMessageSend(eventContainer.MessageCreate.ChannelID, "pong")
				}
			}
		}
	}
}
