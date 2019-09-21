package greeter

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/Cacophony/Processor/plugins/automod/models"
)

type greeterType int

const (
	greeterTypeJoin greeterType = iota
	greeterTypeLeave
)

func (gt greeterType) String() string {
	switch gt {
	case greeterTypeJoin:
		return "Join"
	case greeterTypeLeave:
		return "Leave"
	}

	return "Unknown"
}

type Entry struct {
	gorm.Model
	GuildID   string
	ChannelID string
	Type      greeterType
	Message   string
	RuleID    uint
	Rule      models.Rule
}

func (*Entry) TableName() string {
	return "greeter_entries"
}
