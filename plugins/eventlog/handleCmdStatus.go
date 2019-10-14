package eventlog

import (
	"strings"

	"gitlab.com/Cacophony/go-kit/config"

	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleCmdStatus(event *events.Event) {
	enabled, err := config.GuildGetBool(p.db, event.GuildID, eventlogEnableKey)
	if err != nil && !strings.Contains(err.Error(), "record not found") {
		event.Except(err)
		return
	}

	_, err = event.Respond("eventlog.status.message", "enabled", enabled)
	event.Except(err)
}