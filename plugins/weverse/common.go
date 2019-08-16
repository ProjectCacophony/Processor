package weverse

import (
	"errors"
	"strings"

	"github.com/Seklfreak/geverse"
	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/events"
)

func feedsPerGuildLimit(event *events.Event) int {
	return -1
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

func extractCommunity(communities []geverse.Community, args []string) (*geverse.Community, error) {
	for _, arg := range args {
		for _, community := range communities {

			if strings.EqualFold(community.Name, arg) {
				community := community
				return &community, nil
			}

			if strings.EqualFold(
				strings.Replace(strings.Join(community.Fullname, ""), " ", "", -1),
				strings.Replace(arg, " ", "", -1),
			) {
				community := community
				return &community, nil
			}
		}
	}

	return nil, errors.New("community not found")
}
