package customcommands

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Command struct {
	gorm.Model

	Name          string
	GuildID       string // Could be blank if IsUserCommand is true and it was made in DMs
	UserID        string
	Content       string
	ObjectName    string
	CreatedAt     time.Time
	Triggered     int
	IsUserCommand bool
}

func (*Command) TableName() string {
	return "custom_commands"
}
