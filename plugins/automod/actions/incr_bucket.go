package actions

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"gitlab.com/Cacophony/Processor/plugins/automod/interfaces"
	"gitlab.com/Cacophony/Processor/plugins/automod/models"
	"gitlab.com/Cacophony/go-kit/bucket"
	"gitlab.com/Cacophony/go-kit/events"
)

type IncrBucket struct {
}

func (t IncrBucket) Name() string {
	return "incr_bucket"
}

func (t IncrBucket) Args() int {
	return 4
}

func (t IncrBucket) NewItem(env *models.Env, args []string) (interfaces.ActionItemInterface, error) {
	if len(args) < 4 {
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

	decay, err := time.ParseDuration(args[2])
	if err != nil {
		return nil, err
	}

	amount, err := strconv.Atoi(args[3])
	if err != nil {
		return nil, err
	}
	if amount < 1 {
		return nil, errors.New("amount has to be bigger than 0")
	}

	// TODO: sanitize bucket tag

	return &IncrBucketItem{
		Decay:     decay,
		Amount:    amount,
		Type:      bucketType, // TODO: support setting other types
		TagSuffix: args[0],
	}, nil
}

func (t IncrBucket) Description() string {
	return "automod.actions.incr_bucket"
}

type IncrBucketItem struct {
	Decay     time.Duration
	Amount    int
	TagSuffix string
	Type      events.BucketType
}

func (t *IncrBucketItem) Do(env *models.Env) {
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

	for _, key := range keys {
		// TODO: add support for amount
		valueList, _ := bucket.AddWithValue(
			env.Redis,
			key,
			bucketContent(env),
			t.Decay,
		)

		values := make([]string, len(valueList))

		for i, value := range valueList {
			stringValue, ok := value.Member.(string)
			if !ok {
				continue
			}

			values[i] = stringValue
		}

		// TODO: publish event?
		// TODO: async
		env.Handler.Handle(&events.Event{
			Type: events.CacophonyBucketUpdate,
			BucketUpdate: &events.BucketUpdate{
				Tag:       t.TagSuffix,
				GuildID:   env.GuildID,
				Values:    values,
				Type:      t.Type,
				KeySuffix: key,
			},
		})
	}
}

func bucketTag(guildID, channelID, userID, tag string) string {
	key := "automod:guild:" + guildID
	if channelID != "" {
		key += ":channel:" + channelID
	}
	if userID != "" {
		key += ":user:" + userID
	}
	key += ":" + tag
	return key
}

func bucketContent(env *models.Env) string {
	return env.GuildID + "|" +
		strings.Join(env.ChannelID, ";") + "|" +
		strings.Join(env.UserID, ";") + "|" +
		strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
}

// TODO: allow user channel, or user specific buckets (change bucket tag?)
