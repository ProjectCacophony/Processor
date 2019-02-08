package main

import (
	"gitlab.com/Cacophony/go-kit/logging"
)

type config struct {
	Environment           logging.Environment `envconfig:"ENVIRONMENT" default:"development"`
	AMQPDSN               string              `envconfig:"AMQP_DSN" default:"amqp://guest:guest@localhost:5672/"`
	LoggingDiscordWebhook string              `envconfig:"LOGGING_DISCORD_WEBHOOK"`
}
