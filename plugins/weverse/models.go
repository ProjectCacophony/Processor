package weverse

import (
	"errors"
	"time"

	"github.com/lib/pq"
	"gitlab.com/Cacophony/go-kit/feed"

	"github.com/jinzhu/gorm"
)

type Entry struct {
	gorm.Model
	GuildID           string
	ChannelOrUserID   string // UserID in case of DMs
	DM                bool
	DisableArtistFeed bool `gorm:"default:false"`
	DisableMediaFeed  bool `gorm:"default:false"`
	DisableNoticeFeed bool `gorm:"default:false"`
	DisableMomentFeed bool `gorm:"default:false"`
	AddedBy           string
	BotID             string // only relevant for DMs

	WeverseChannelName string
	WeverseChannelID   int64

	ArtistFeedLastCheck time.Time
	ArtistFeedCheck     feed.Check `gorm:"embedded;embedded_prefix:artist_"`

	MediaFeedLastCheck time.Time
	MediaFeedCheck     feed.Check `gorm:"embedded;embedded_prefix:media_"`

	NoticeFeedLastCheck time.Time
	NoticeFeedCheck     feed.Check `gorm:"embedded;embedded_prefix:notice_"`

	MomentFeedLastCheck time.Time
	MomentFeedCheck     feed.Check `gorm:"embedded;embedded_prefix:moment_"`
}

func (*Entry) TableName() string {
	return "weverse_entries"
}

// Post model maintained by Worker
type Post struct {
	gorm.Model
	EntryID uint

	PostType   string
	PostID     int64
	MessageIDs pq.StringArray `gorm:"type:varchar[]"`
}

func (*Post) TableName() string {
	return "weverse_posts"
}

type modifyType int

const (
	modifyArtist modifyType = iota
	modifyMedia
	modifyNotice
	modifyMoment
)

func entryModify(db *gorm.DB, id uint, modification modifyType, value bool) error {
	var fieldName string
	switch modification {
	case modifyArtist:
		fieldName = "disable_artist_feed"
	case modifyMedia:
		fieldName = "disable_media_feed"
	case modifyNotice:
		fieldName = "disable_notice_feed"
	case modifyMoment:
		fieldName = "disable_moment_feed"
	}

	if fieldName == "" {
		return errors.New("invalid modification type")
	}

	return db.Model(&Entry{}).Where("id = ?", id).Update(fieldName, value).Error
}
