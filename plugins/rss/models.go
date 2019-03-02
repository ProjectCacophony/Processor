package rss

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
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

// Post model maintained by Worker
type Post struct {
	gorm.Model
	EntryID uint

	PostGUID   string
	PostLink   string
	MessageIDs pq.StringArray `gorm:"type:varchar[]"`
}

func (*Post) TableName() string {
	return "rss_posts"
}
