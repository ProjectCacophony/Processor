package instagram

import (
	"time"

	"gitlab.com/Cacophony/go-kit/feed"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

type Entry struct {
	gorm.Model
	GuildID          string
	ChannelOrUserID  string // UserID in case of DMs
	DM               bool
	DisablePostFeed  bool `gorm:"default:false"`
	DisableStoryFeed bool `gorm:"default:false"`
	AddedBy          string
	BotID            string // only relevant for DMs

	InstagramUsername  string
	InstagramAccountID string

	LastCheck time.Time
	feed.Check

	StoriesLastCheck time.Time
	StoriesCheck     feed.Check `gorm:"embedded;embedded_prefix:stories_"`

	IGTVLastCheck time.Time
	IGTVCheck     feed.Check `gorm:"embedded;embedded_prefix:igtv_"`
}

func (*Entry) TableName() string {
	return "instagram_entries"
}

// Post model maintained by Worker
type Post struct {
	gorm.Model
	EntryID uint

	PostID     string
	MessageIDs pq.StringArray `gorm:"type:varchar[]"`
}

func (*Post) TableName() string {
	return "instagram_posts"
}
