package filters

import (
	"errors"
	"strings"

	"gitlab.com/Cacophony/Processor/plugins/automod/models"

	"gitlab.com/Cacophony/Processor/plugins/automod/interfaces"
	"gitlab.com/Cacophony/go-kit/events"
)

type Channel struct {
}

func (f Channel) Name() string {
	return "if_channel"
}

func (f Channel) Args() int {
	return 1
}

func (f Channel) Deprecated() bool {
	return false
}

func (f Channel) NewItem(env *models.Env, args []string) (interfaces.FilterItemInterface, error) {
	if len(args) < 1 {
		return nil, errors.New("too few arguments")
	}

	channelMentions := strings.Split(args[0], ",")
	channelIDs := make([]string, 0, len(channelMentions))

	for _, channelMention := range channelMentions {
		channel, err := env.State.ChannelFromMention(env.GuildID, channelMention)
		if err != nil {
			continue
		}

		channelIDs = append(channelIDs, channel.ID)
	}

	if len(channelIDs) == 0 {
		return nil, errors.New("no Channel IDs defined")
	}

	return &ChannelItem{
		channelIDs: channelIDs,
	}, nil
}

func (f Channel) Description() string {
	return "automod.filters.if_channel"
}

type ChannelItem struct {
	channelIDs []string
}

func (f *ChannelItem) Match(env *models.Env) bool {
	if env.Event == nil {
		return false
	}

	if env.Event.Type != events.MessageCreateType {
		return false
	}

	for _, envChannelID := range env.ChannelID {
		if matchChannels(f.channelIDs, envChannelID) {
			continue
		}

		return false
	}

	return true
}
