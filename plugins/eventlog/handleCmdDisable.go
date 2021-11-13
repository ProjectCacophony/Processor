package eventlog

import (
	"strings"

	"gitlab.com/Cacophony/go-kit/config"
	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleCmdDisable(event *events.Event) {
	enabled, err := config.GuildGetBool(p.db, event.GuildID, eventlogEnableKey)
	if err != nil && !strings.Contains(err.Error(), "record not found") {
		event.Except(err)
		return
	}

	if !enabled {
		event.Respond("eventlog.disable.already-disabled")
		return
	}

	err = config.GuildSetBool(p.db, event.GuildID, eventlogEnableKey, false)
	if err != nil {
		event.Except(err)
		return
	}
	err = config.GuildSetString(p.db, event.GuildID, eventlogChannelKey, "")
	if err != nil {
		event.Except(err)
		return
	}

	_, err = event.Respond("eventlog.disable.success")
	event.Except(err)
}
