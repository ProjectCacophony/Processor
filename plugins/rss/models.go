package rss

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Entry struct {
	gorm.Model
	GuildID   string
	ChannelID string
	AddedBy   string

	Name    string
	URL     string
	FeedURL string

	LastCheck time.Time
}

func (*Entry) TableName() string {
	return "rss_entries"
}
