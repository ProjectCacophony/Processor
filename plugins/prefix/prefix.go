package prefix

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/Cacophony/go-kit/config"
	"gitlab.com/Cacophony/go-kit/events"
)

func handleGetPrefix(event *events.Event) {
	// event.Prefix() is automatically called and will query for guild prefix
	_, err := event.Respond("prefix.get-prefix")
	event.Except(err)
}

func handleSetPrefix(event *events.Event, db *gorm.DB) {
	err := event.Discord().Client.ChannelTyping(event.ChannelID)
	event.Except(err)

	if len(event.Fields()) != 3 {
		_, err := event.Respond("prefix.set-prefix.no-value")
		event.Except(err)
		return
	}

	newPrefix := event.Fields()[2]
	err = config.GuildSetString(db, event.GuildID, guildCmdPrefixKey, newPrefix)
	if err != nil {
		event.Except(err)
		return
	}

	_, err = event.Respond("prefix.set-prefix.success", "newPrefix", newPrefix)
	event.Except(err)
}
