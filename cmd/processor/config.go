package main

import (
	"time"

	"gitlab.com/Cacophony/go-kit/logging"
)

// nolint: lll
type config struct {
	Port                      int                 `envconfig:"PORT" default:"8000"`
	Environment               logging.Environment `envconfig:"ENVIRONMENT" default:"development"`
	AMQPDSN                   string              `envconfig:"AMQP_DSN" default:"amqp://guest:guest@localhost:5672/"`
	LoggingDiscordWebhook     string              `envconfig:"LOGGING_DISCORD_WEBHOOK"`
	ConcurrentProcessingLimit int                 `envconfig:"CONCURRENT_PROCESSING_LIMIT" default:"50"`
	ProcessingDeadline        time.Duration       `envconfig:"PROCESSING_DEADLINE" default:"5m"`
	DiscordTokens             map[string]string   `envconfig:"DISCORD_TOKENS"`
	DBDSN                     string              `envconfig:"DB_DSN" default:"postgres://postgres:postgres@localhost:5432/?sslmode=disable"`
	RedisAddress              string              `envconfig:"REDIS_ADDRESS" default:"localhost:6379"`
	RedisPassword             string              `envconfig:"REDIS_PASSWORD"`
}
