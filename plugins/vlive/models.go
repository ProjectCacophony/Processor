package vlive

import (
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

	VLiveChannelID string
}

func (*Entry) TableName() string {
	return "vlive_entries"
}

// Post model maintained by Worker
type Post struct {
	gorm.Model
	EntryID uint

	PostID     string
	MessageIDs pq.StringArray `gorm:"type:varchar[]"`
}

func (*Post) TableName() string {
	return "vlive_posts"
}
