package actions

import (
	"errors"

	"gitlab.com/Cacophony/Processor/plugins/automod/interfaces"
	"gitlab.com/Cacophony/Processor/plugins/automod/models"
	"gitlab.com/Cacophony/go-kit/bucket"
	"gitlab.com/Cacophony/go-kit/events"
)

type ResetBucket struct{}

func (t ResetBucket) Name() string {
	return "reset_bucket"
}

func (t ResetBucket) Args() int {
	return 2
}

func (t ResetBucket) Deprecated() bool {
	return false
}

func (t ResetBucket) NewItem(env *models.Env, args []string) (interfaces.ActionItemInterface, error) {
	if len(args) < 2 {
		return nil, errors.New("too few arguments")
	}

	var bucketType events.BucketType
	switch args[1] {
	case "guild":
		bucketType = events.GuildBucketType
	case "channel":
		bucketType = events.ChannelBucketType
	case "user":
		bucketType = events.UserBucketType
	default:
		return nil, errors.New("invalid bucket type")
	}

	// TODO: sanitize bucket tag

	return &ResetBucketItem{
		Type:      bucketType,
		TagSuffix: args[0],
	}, nil
}

func (t ResetBucket) Description() string {
	return "automod.actions.reset_bucket"
}

type ResetBucketItem struct {
	TagSuffix string
	Type      events.BucketType
}

func (t *ResetBucketItem) Do(env *models.Env) (bool, error) {
	var keys []string
	switch t.Type {
	case events.GuildBucketType:
		keys = []string{bucketTag(env.GuildID, "", "", t.TagSuffix)}
	case events.ChannelBucketType:
		for _, channelID := range env.ChannelID {
			keys = []string{bucketTag(env.GuildID, channelID, "", t.TagSuffix)}
		}
	case events.UserBucketType:
		for _, userID := range env.UserID {
			keys = []string{bucketTag(env.GuildID, "", userID, t.TagSuffix)}
		}
	}

	var err error

	for _, key := range keys {

		err = bucket.Reset(
			env.Redis,
			key,
		)
		if err != nil {
			return false, err
		}
	}

	return false, nil
}
