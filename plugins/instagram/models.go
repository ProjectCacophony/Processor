package instagram

import (
	"time"

	"gitlab.com/Cacophony/go-kit/feed"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

type Entry struct {
	gorm.Model
	GuildID         string
	ChannelOrUserID string // UserID in case of DMs
	DM              bool
	AddedBy         string
	BotID           string // only relevant for DMs

	InstagramUsername  string
	InstagramAccountID string

	LastCheck time.Time
	feed.Check
}

func (*Entry) TableName() string {
	return "instagram_entries"
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
	return "instagram_posts"
}
