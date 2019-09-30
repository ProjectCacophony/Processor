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

func (t BucketUpdated) Deprecated() bool {
	return false
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
		userIDs, channelIDs, GuildID, messages := extractBucketValues(value)

		env.GuildID = GuildID

		for _, userID := range userIDs {
			if userID == "" {
				continue
			}

			if !stringSliceContains(env.UserID, userID) {
				env.UserID = append(env.UserID, userID)
			}
		}

		for _, channelID := range channelIDs {
			if channelID == "" {
				continue
			}

			if !stringSliceContains(env.ChannelID, channelID) {
				env.ChannelID = append(env.ChannelID, channelID)
			}
		}

		for _, message := range messages {
			if message.ID == "" || message.ChanneID == "" {
				continue
			}

			if !messageSliceContains(env.Messages, message) {
				env.Messages = append(env.Messages, message)
			}
		}
	}

	return true
}

func stringSliceContains(haystack []string, needle string) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}

	return false
}

func messageSliceContains(haystack []*models.EnvMessage, needle *models.EnvMessage) bool {
	for _, item := range haystack {
		if item.ID == needle.ID && item.ChanneID == needle.ChanneID {
			return true
		}
	}

	return false
}

func extractBucketValues(value string) (userIDs, channelIDs []string, guildID string, messages []*models.EnvMessage) {
	parts := strings.Split(value, "|")
	if len(parts) < 4 {
		return
	}

	guildID = parts[0]

	channelIDs = strings.Split(parts[1], ";")

	userIDs = strings.Split(parts[2], ";")

	for _, messageData := range strings.Split(parts[3], ";") {
		messageParts := strings.Split(messageData, ":")
		if len(messageParts) < 2 {
			continue
		}

		messages = append(messages, &models.EnvMessage{
			ID:       messageParts[0],
			ChanneID: messageParts[1],
		})
	}

	return
}
