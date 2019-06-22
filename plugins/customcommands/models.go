package customcommands

import (
	"errors"

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
	Triggered     int
	IsUserCommand bool
}

func (*CustomCommand) TableName() string {
	return "custom_commands"
}

func (c *CustomCommand) getContent() string {
	output := c.Content

	if c.File != nil && c.File.GetLink() != "" {
		if output != "" {
			output += "\n"
		}
		output += c.File.GetLink()
	}
	return output
}

func (c *CustomCommand) run(event *events.Event) error {
	err := c.triggered(event.DB())
	if err != nil {
		return err
	}

	_, err = event.Respond(c.getContent())
	return err
}

func (c *CustomCommand) triggered(db *gorm.DB) error {
	if c == nil {
		return errors.New("command cannot be nil")
	}
	return db.Model(c).Update("triggered", gorm.Expr("triggered + 1")).Error
}
