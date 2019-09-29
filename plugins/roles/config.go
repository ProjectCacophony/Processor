package roles

import (
	"gitlab.com/Cacophony/go-kit/config"
	"gitlab.com/Cacophony/go-kit/events"
)

const guildRoleChannelKey = "cacophony:processor:role:default-role-channel"

func (p *Plugin) setRoleChannel(event *events.Event) {
	if len(event.Fields()) < 3 {
		event.Respond("common.invalid-params")
		return
	}

	channel, err := event.State().ChannelFromMention(event.GuildID, event.Fields()[2])
	if err != nil {
		event.Except(err)
		return
	}

	err = config.GuildSetString(event.DB(), event.GuildID, guildRoleChannelKey, channel.ID)
	if err != nil {
		event.Except(err)
		return
	}

	_, err = event.Respond("roles.config.channel", "channelName", channel.Name)
	if err != nil {
		event.Except(err)
		return
	}
}
