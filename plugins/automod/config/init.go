package config

import (
	"github.com/jinzhu/gorm"
)

// put package in go-kit

func InitConfig(db *gorm.DB) error {
	return db.AutoMigrate(Item{}).Error
}
