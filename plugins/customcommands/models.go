package customcommands

import (
	"errors"
	"strings"

	"github.com/jinzhu/gorm"
	"gitlab.com/Cacophony/go-kit/discord"
	"gitlab.com/Cacophony/go-kit/events"
	"go.uber.org/zap"
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
	Type          customCommandType
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

	switch c.Type {
	case customCommandTypeContent:
		message := discord.MessageCodeToMessage(c.getContent())
		_, err = event.RespondComplex(message)
		return err
	case customCommandTypeCommand:
		event.MessageCreate.Content = event.Prefix() +
			strings.TrimSpace(c.getContent()+" "+strings.Join(event.Fields()[1:], " "))

		event.Logger().Info("executing custom command",
			zap.Uint("customcommand_id", c.ID),
			zap.String("final_command", event.MessageCreate.Content),
		)

		err, recoverable := event.Publisher().Publish(event.Context(), event)
		if err != nil {
			if !recoverable {
				event.Logger().Fatal(
					"received unrecoverable error while publishing custom commands alias message",
					zap.Error(err),
				)
			}
		}
		return err
	}

	return nil
}

func (c *CustomCommand) triggered(db *gorm.DB) error {
	if c == nil {
		return errors.New("command cannot be nil")
	}
	return db.Model(c).Update("triggered", gorm.Expr("triggered + 1")).Error
}

type customCommandType int

const (
	customCommandTypeContent customCommandType = iota
	customCommandTypeCommand
)
