package eventlog

import (
	"strings"

	"gitlab.com/Cacophony/go-kit/config"

	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleCmdEnable(event *events.Event) {
	enabled, err := config.GuildGetBool(p.db, event.GuildID, eventlogEnableKey)
	if err != nil && !strings.Contains(err.Error(), "record not found") {
		event.Except(err)
		return
	}

	targetChannel, err := event.FindChannel()
	if err != nil {
		event.Except(err)
		return
	}

	eventlogChannelID, err := config.GuildGetString(p.db, event.GuildID, eventlogChannelKey)
	if err != nil && !strings.Contains(err.Error(), "record not found") {
		event.Except(err)
		return
	}

	if enabled && eventlogChannelID == targetChannel.ID {
		event.Respond("eventlog.enable.already-enabled")
		return
	}

	if !enabled {
		err = config.GuildSetBool(p.db, event.GuildID, eventlogEnableKey, true)
		if err != nil {
			event.Except(err)
			return
		}
	}
	if eventlogChannelID != targetChannel.ID {
		err = config.GuildSetString(p.db, event.GuildID, eventlogChannelKey, targetChannel.ID)
		if err != nil {
			event.Except(err)
			return
		}
	}

	_, err = event.Respond("eventlog.enable.success", "channel", targetChannel)
	event.Except(err)
}
