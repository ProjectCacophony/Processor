package automod

import (
	"github.com/jinzhu/gorm"
)

type Rule struct {
	gorm.Model
	GuildID string `gorm:"index;not null;unique_index:idx_automod_rules_guildid_name"`
	Name    string `gorm:"not null;unique_index:idx_automod_rules_guildid_name"`
	Trigger string `gorm:"not null"`
	Filters []RuleFilter
	Actions []RuleAction
	Process bool
}

func (*Rule) TableName() string {
	return "automod_rules"
}

type RuleFilter struct {
	gorm.Model
	RuleID uint
	Name   string `gorm:"not null"`
	Value  string
}

func (*RuleFilter) TableName() string {
	return "automod_rule_filters"
}

type RuleAction struct {
	gorm.Model
	RuleID uint
	Name   string `gorm:"not null"`
	Value  string
}

func (*RuleAction) TableName() string {
	return "automod_rule_actions"
}
