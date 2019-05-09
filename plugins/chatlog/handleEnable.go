package chatlog

import (
	"strings"

	"gitlab.com/Cacophony/go-kit/config"

	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleEnable(event *events.Event) {
	enabled, err := config.GuildGetBool(p.db, event.GuildID, chatlogEnabledKey)
	if err != nil && !strings.Contains(err.Error(), "record not found") {
		event.Except(err)
		return
	}

	if enabled {
		event.Respond("chatlog.enable.already-enabled")
		return
	}

	err = config.GuildSetBool(p.db, event.GuildID, chatlogEnabledKey, true)
	if err != nil {
		event.Except(err)
		return
	}

	_, err = event.Respond("chatlog.enable.success")
	event.Except(err)
}
