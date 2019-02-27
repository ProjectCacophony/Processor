package gall

import (
	"strings"

	"gitlab.com/Cacophony/go-kit/permissions"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/events"
)

func feedsPerGuildLimit(event *events.Event) int {
	if event.Has(permissions.BotOwner) {
		return -1
	}

	return 2
}

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

func paramsIsAll(args []string) (bool, []string) {
	for i, arg := range args {
		if !strings.EqualFold(arg, "all") {
			continue
		}

		return true, append(args[:i], args[i+1:]...)
	}
	return false, args
}
