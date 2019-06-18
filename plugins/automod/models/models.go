package models

import (
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

type Rule struct {
	gorm.Model
	GuildID       string         `gorm:"index;not null"`
	Name          string         `gorm:"not null"`
	TriggerName   string         `gorm:"not null"`
	TriggerValues pq.StringArray `gorm:"type:varchar[]"`
	Filters       []RuleFilter
	Actions       []RuleAction
	Stop          bool
	Silent        bool
	Runs          int `gorm:"default:'0'"`
}

func (*Rule) TableName() string {
	return "automod_rules"
}

type RuleFilter struct {
	gorm.Model
	RuleID uint
	Name   string         `gorm:"not null"`
	Values pq.StringArray `gorm:"type:varchar[]"`
	Not    bool
}

func (*RuleFilter) TableName() string {
	return "automod_rule_filters"
}

type RuleAction struct {
	gorm.Model
	RuleID uint
	Name   string         `gorm:"not null"`
	Values pq.StringArray `gorm:"type:varchar[]"`
}

func (*RuleAction) TableName() string {
	return "automod_rule_actions"
}

type LogEntry struct {
	gorm.Model
	GuildID      string
	Rule         Rule
	RuleID       uint
	ChannelIDs   pq.StringArray `gorm:"type:varchar[]"`
	UserIDs      pq.StringArray `gorm:"type:varchar[]"`
	MessageIDs   pq.StringArray `gorm:"type:varchar[]"`
	ErrorMessage string
}

func (*LogEntry) TableName() string {
	return "automod_log"
}
