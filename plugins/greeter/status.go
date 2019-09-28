package greeter

import (
	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleStatus(event *events.Event) {
	entries, err := entriesFind(event.DB(), event.GuildID)
	if err != nil {
		event.Except(err)
		return
	}

	currentChannel, err := event.State().Channel(event.ChannelID)
	if err != nil {
		event.Except(err)
		return
	}

	event.Respond("greeter.status.content", "entries", entries, "channel", currentChannel)
}
