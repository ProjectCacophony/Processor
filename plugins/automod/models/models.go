package models

import (
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
)

type Rule struct {
	gorm.Model
	GuildID       string         `gorm:"index;not null;unique_index:idx_automod_rules_guildid_name"`
	Name          string         `gorm:"not null;unique_index:idx_automod_rules_guildid_name"`
	TriggerName   string         `gorm:"not null"`
	TriggerValues pq.StringArray `gorm:"type:varchar[]"`
	Filters       []RuleFilter
	Actions       []RuleAction
	Process       bool
}

func (*Rule) TableName() string {
	return "automod_rules"
}

type RuleFilter struct {
	gorm.Model
	RuleID uint
	Name   string         `gorm:"not null"`
	Values pq.StringArray `gorm:"type:varchar[]"`
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
