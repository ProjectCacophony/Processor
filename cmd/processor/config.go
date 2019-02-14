package main

import (
	"gitlab.com/Cacophony/go-kit/logging"
)

type config struct {
	Port                  int                 `envconfig:"PORT" default:"8000"`
	Environment           logging.Environment `envconfig:"ENVIRONMENT" default:"development"`
	AMQPDSN               string              `envconfig:"AMQP_DSN" default:"amqp://guest:guest@localhost:5672/"`
	LoggingDiscordWebhook string              `envconfig:"LOGGING_DISCORD_WEBHOOK"`
}