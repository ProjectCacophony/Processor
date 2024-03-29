package gall

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/permissions"
)

func feedsPerGuildLimit(event *events.Event) int {
	if event.Has(permissions.BotAdmin) {
		return -1
	}

	if event.DM() {
		return 2
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
