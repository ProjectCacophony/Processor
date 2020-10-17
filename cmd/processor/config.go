package main

import (
	"time"

	"gitlab.com/Cacophony/go-kit/errortracking"
	"gitlab.com/Cacophony/go-kit/featureflag"
	"gitlab.com/Cacophony/go-kit/logging"
)

type config struct {
	Port                      int                  `envconfig:"PORT" default:"8000"`
	Hash                      string               `envconfig:"HASH"`
	Environment               logging.Environment  `envconfig:"ENVIRONMENT" default:"development"`
	ClusterEnvironment        string               `envconfig:"CLUSTER_ENVIRONMENT" default:"development"`
	AMQPDSN                   string               `envconfig:"AMQP_DSN" default:"amqp://guest:guest@localhost:5672/"`
	LoggingDiscordWebhook     string               `envconfig:"LOGGING_DISCORD_WEBHOOK"`
	ConcurrentProcessingLimit int                  `envconfig:"CONCURRENT_PROCESSING_LIMIT" default:"50"`
	ProcessingDeadline        time.Duration        `envconfig:"PROCESSING_DEADLINE" default:"5m"`
	DiscordTokens             map[string]string    `envconfig:"DISCORD_TOKENS"`
	DBDSN                     string               `envconfig:"DB_DSN" default:"postgres://postgres:postgres@localhost:5432/?sslmode=disable"`
	RedisAddress              string               `envconfig:"REDIS_ADDRESS" default:"localhost:6379"`
	RedisPassword             string               `envconfig:"REDIS_PASSWORD"`
	FeatureFlag               featureflag.Config   `envconfig:"FEATUREFLAG"`
	ErrorTracking             errortracking.Config `envconfig:"ERRORTRACKING"`
	DiscordAPIBase            string               `envconfig:"DISCORD_API_BASE"`
	GoogleCloudBucketName     string               `envconfig:"GCLOUD_BUCKET_NAME"`
	ObjectStorageFQDN         string               `envconfig:"OBJECT_STORAGE_FQDN"`
	HoneycombAPIKey           string               `envconfig:"HONEYCOMB_API_KEY"`
}
