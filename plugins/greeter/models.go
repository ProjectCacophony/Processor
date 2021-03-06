package greeter

import (
	"time"

	"github.com/jinzhu/gorm"
	"gitlab.com/Cacophony/Processor/plugins/automod/models"
)

type greeterType int

const (
	greeterTypeJoin greeterType = iota
	greeterTypeLeave
	greeterTypeBan
	greeterTypeUnban
)

func (gt greeterType) String() string {
	switch gt {
	case greeterTypeJoin:
		return "Join"
	case greeterTypeLeave:
		return "Leave"
	case greeterTypeBan:
		return "Ban"
	case greeterTypeUnban:
		return "Unban"
	}

	return "Unknown"
}

type Entry struct {
	gorm.Model
	GuildID    string
	ChannelID  string
	Type       greeterType
	Message    string
	AutoDelete time.Duration
	RuleID     uint
	Rule       models.Rule
}

func (*Entry) TableName() string {
	return "greeter_entries"
}
