package automod

import (
	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/events"
)

func paramsExtractChannel(event *events.Event, args []string) (*discordgo.Channel, []string, error) {
	for i, arg := range args {
		channel, err := event.State().ChannelFromMention(event.GuildID, arg)
		if err != nil {
			continue
		}

		return channel, append(args[:i], args[i+1:]...), nil
	}

	channel, err := event.State().Channel(event.ChannelID)
	return channel, args, err
}
