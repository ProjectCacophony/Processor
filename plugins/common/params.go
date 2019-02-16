package common

import "github.com/jinzhu/gorm"

type StartParameters struct {
	DB *gorm.DB
}

type StopParameters struct {
	DB *gorm.DB
}
