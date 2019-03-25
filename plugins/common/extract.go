package common

import (
	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/events"
)

func FieldsExtractChannelTypes(
	event *events.Event, args []string, types ...discordgo.ChannelType,
) (*discordgo.Channel, []string, error) {
	for i, arg := range args {
		channel, err := event.State().ChannelFromMentionTypes(event.GuildID, arg, types...)
		if err != nil {
			continue
		}

		return channel, append(args[:i], args[i+1:]...), nil
	}

	channel, err := event.State().Channel(event.ChannelID)
	return channel, args, err
}
