package chatlog

import (
	"strings"

	"gitlab.com/Cacophony/go-kit/config"

	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleDisable(event *events.Event) {
	enabled, err := config.GuildGetBool(p.db, event.GuildID, chatlogEnabledKey)
	if err != nil && !strings.Contains(err.Error(), "record not found") {
		event.Except(err)
		return
	}

	if !enabled {
		event.Respond("chatlog.disable.already-disabled") // nolint: errcheck
		return
	}

	err = config.GuildSetBool(p.db, event.GuildID, chatlogEnabledKey, false)
	if err != nil {
		event.Except(err)
		return
	}

	_, err = event.Respond("chatlog.disable.success")
	event.Except(err)
}
