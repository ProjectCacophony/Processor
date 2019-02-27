package gall

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Entry struct {
	gorm.Model
	GuildID   string
	ChannelID string
	AddedBy   string

	BoardID      string
	MinorGallery bool
	Recommended  bool
	LastCheck    time.Time
}

func (*Entry) TableName() string {
	return "gall_entries"
}
