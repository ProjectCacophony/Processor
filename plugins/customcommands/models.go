package customcommands

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
	"gitlab.com/Cacophony/go-kit/events"
)

type CustomCommand struct {
	gorm.Model

	Name          string
	GuildID       string // Could be blank if IsUserCommand is true and it was made in DMs
	UserID        string
	Content       string
	File          *events.FileInfo `gorm:"foreignkey:CustomCommandID"`
	Date          time.Time
	Triggered     int
	IsUserCommand bool
}

func (*CustomCommand) TableName() string {
	return "custom_commands"
}

func (c *CustomCommand) run(event *events.Event) error {
	err := c.triggered(event.DB())
	if err != nil {
		return err
	}

	output := c.Content
	if c.File != nil {
		output += "\n" + c.File.GetLink()
	}

	_, err = event.Respond(output)
	return err
}

func (c *CustomCommand) triggered(db *gorm.DB) error {
	if c == nil {
		return errors.New("command cannot be nil")
	}
	return db.Model(c).Update("triggered", gorm.Expr("triggered + 1")).Error
}
