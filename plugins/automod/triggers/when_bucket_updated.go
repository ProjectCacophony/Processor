package triggers

import (
	"errors"
	"strings"

	"gitlab.com/Cacophony/Processor/plugins/automod/interfaces"
	"gitlab.com/Cacophony/Processor/plugins/automod/models"
	"gitlab.com/Cacophony/go-kit/events"
)

type BucketUpdated struct {
}

func (t BucketUpdated) Name() string {
	return "when_bucket_updated"
}

func (t BucketUpdated) Args() int {
	return 1
}

func (t BucketUpdated) NewItem(env *models.Env, args []string) (interfaces.TriggerItemInterface, error) {
	if len(args) < 1 {
		return nil, errors.New("too few arguments")
	}

	return &BucketUpdatedItem{
		Tag: args[0],
	}, nil
}

func (t BucketUpdated) Description() string {
	return "automod.triggers.when_bucket_updated"
}

type BucketUpdatedItem struct {
	Tag string
}

func (t *BucketUpdatedItem) Match(env *models.Env) bool {
	if env.Event == nil {
		return false
	}

	if env.Event.Type != events.CacophonyBucketUpdate {
		return false
	}

	if env.Event.BucketUpdate.Tag != t.Tag {
		return false
	}

	env.GuildID = env.Event.BucketUpdate.GuildID

	for _, value := range env.Event.BucketUpdate.Values {
		userIDs, channelIDs, GuildID := extractBucketValues(value)
		env.GuildID = GuildID
		env.ChannelID = append(env.ChannelID, channelIDs...)
		env.UserID = append(env.UserID, userIDs...)
	}

	return true
}

func extractBucketValues(value string) (userIDs, channelIDs []string, guildID string) {
	parts := strings.Split(value, "|")
	if len(parts) < 3 {
		return
	}

	guildID = parts[0]

	channelIDs = strings.Split(parts[1], ";")

	userIDs = strings.Split(parts[2], ";")

	return
}
