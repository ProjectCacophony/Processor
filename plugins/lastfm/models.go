package lastfm

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	UserID         string `gorm:"unique_index"`
	LastFMUsername string
}

func (*User) TableName() string {
	return "lastfm_users"
}
