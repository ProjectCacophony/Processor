package chatlog

import (
	"strings"

	"gitlab.com/Cacophony/go-kit/config"

	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleStatus(event *events.Event) {
	enabled, err := config.GuildGetBool(p.db, event.GuildID, chatlogEnabledKey)
	if err != nil && !strings.Contains(err.Error(), "record not found") {
		event.Except(err)
		return
	}

	_, err = event.Respond("chatlog.status.message", "enabled", enabled)
	event.Except(err)
}
