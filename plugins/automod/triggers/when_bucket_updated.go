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

	var newEnv *models.Env
	var parts []string
	for _, envData := range env.Event.BucketUpdate.EnvDatas {
		parts = strings.Split(string(envData), "|")
		if len(parts) < 2 {
			continue
		}
		envData = []byte(parts[1])

		newEnv = &models.Env{}
		err := newEnv.Unmarshal(envData)
		if err != nil {
			newEnv.Event.ExceptSilent(err)
		}

		env.GuildID = newEnv.GuildID
		for _, userID := range newEnv.UserID {
			if userID == "" {
				continue
			}

			if !stringSliceContains(env.UserID, userID) {
				env.UserID = append(env.UserID, userID)
			}
		}

		for _, channelID := range newEnv.ChannelID {
			if channelID == "" {
				continue
			}

			if !stringSliceContains(env.ChannelID, channelID) {
				env.ChannelID = append(env.ChannelID, channelID)
			}
		}

		for _, message := range newEnv.Messages {
			if message.ID == "" || message.ChannelID == "" {
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
		if item.ID == needle.ID && item.ChannelID == needle.ChannelID {
			return true
		}
	}

	return false
}
