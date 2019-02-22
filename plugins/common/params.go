package common

import (
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

type StartParameters struct {
	Logger *zap.Logger
	DB     *gorm.DB
	Redis  *redis.Client
	Tokens map[string]string
}

type StopParameters struct {
	Logger *zap.Logger
	DB     *gorm.DB
	Redis  *redis.Client
	Tokens map[string]string
}
