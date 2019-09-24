package roles

import "github.com/jinzhu/gorm"

type Category struct {
	gorm.Model

	Name      string
	GuildID   string
	ChannelID string // the channel this Category will listen to for role assignments
	Message   string
	Pool      string
	Roles     []Role
	Limit     int
	Enabled   bool
}

type Role struct {
	gorm.Model

	Name      string
	PrintName string
	Enabled   bool
	Aliases   []string
	// Reactions []string
}

func (*Category) TableName() string {
	return "role_categories"
}

func (*Role) TableName() string {
	return "roles"
}
