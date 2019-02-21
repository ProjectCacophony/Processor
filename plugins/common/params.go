package common

import (
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

type StartParameters struct {
	Logger *zap.Logger
	DB     *gorm.DB
}

type StopParameters struct {
	Logger *zap.Logger
	DB     *gorm.DB
}
