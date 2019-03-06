package automod

import (
	"gitlab.com/Cacophony/Processor/plugins/automod/config"
	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) cmdLog(event *events.Event) {
	channel, _, err := paramsExtractChannel(event, event.Fields())
	if err != nil {
		event.Except(err)
		return
	}

	err = config.GuildSetString(p.db, event.GuildID, automodLogKey, channel.ID)
	if err != nil {
		event.Except(err)
		return
	}

	_, err = event.Respond("automod.log.response", "channel", channel)
	event.Except(err)
}
