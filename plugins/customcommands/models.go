package customcommands

import (
	"time"

	"github.com/jinzhu/gorm"
	"gitlab.com/Cacophony/go-kit/events"
)

type Entry struct {
	gorm.Model

	Name          string
	GuildID       string // Could be blank if IsUserCommand is true and it was made in DMs
	UserID        string
	Content       string
	ObjectName    string
	Date          time.Time
	Triggered     int
	IsUserCommand bool
}

func (*Entry) TableName() string {
	return "custom_commands"
}

func (e *Entry) run(event *events.Event) error {
	err := entryUpdateTriggered(event.DB(), e)
	if err != nil {
		return err
	}
	_, err = event.Respond(e.Content)
	return err
}
