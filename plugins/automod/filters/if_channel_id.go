package filters

import (
	"errors"
	"strings"

	"gitlab.com/Cacophony/Processor/plugins/automod/models"

	"gitlab.com/Cacophony/Processor/plugins/automod/interfaces"
	"gitlab.com/Cacophony/go-kit/events"
)

type ChannelID struct {
}

func (f ChannelID) Name() string {
	return "if_channel_id"
}

func (f ChannelID) Args() int {
	return 1
}

func (f ChannelID) NewItem(env *models.Env, args []string) (interfaces.FilterItemInterface, error) {
	if len(args) < 1 {
		return nil, errors.New("too few arguments")
	}

	channelIDs := strings.Split(args[0], ",")

	if len(channelIDs) == 0 {
		return nil, errors.New("no Channel IDs defined")
	}

	return &ChannelIDItem{
		channelIDs: channelIDs,
	}, nil
}

func (f ChannelID) Description() string {
	return "automod.filters.if_channel_id"
}

type ChannelIDItem struct {
	channelIDs []string
}

func (f *ChannelIDItem) Match(env *models.Env) bool {
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

func matchChannels(channels []string, channel string) bool {
	for _, c := range channels {
		if c != channel {
			continue
		}

		return true
	}

	return false
}
