package actions

import (
	"context"
	"errors"
	"math/rand"
	"strconv"
	"time"

	"gitlab.com/Cacophony/Processor/plugins/automod/interfaces"
	"gitlab.com/Cacophony/Processor/plugins/automod/models"
	"gitlab.com/Cacophony/go-kit/bucket"
	"gitlab.com/Cacophony/go-kit/events"
	"go.uber.org/zap"
)

type IncrBucket struct{}

func (t IncrBucket) Name() string {
	return "incr_bucket"
}

func (t IncrBucket) Args() int {
	return 4
}

func (t IncrBucket) Deprecated() bool {
	return false
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
		Type:      bucketType,
		TagSuffix: args[0],
		Random:    rand.New(rand.NewSource(time.Now().UnixNano())),
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
	Random    *rand.Rand
}

func (t *IncrBucketItem) Do(env *models.Env) (bool, error) {
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
	var bucketContentValue []byte
	var values [][]byte
	var recoverable bool

	for _, key := range keys {

		for i := 0; i < t.Amount; i++ {
			bucketContentValue, err = bucketContent(env)
			if err != nil {
				return false, err
			}

			valueList, _ := bucket.AddWithValue(
				env.Redis,
				key,
				bucketContentValue,
				t.Decay,
			)

			values = make([][]byte, len(valueList))

			for i, value := range valueList {
				stringValue, ok := value.Member.(string)
				if !ok {
					continue
				}

				values[i] = []byte(stringValue)
			}
		}

		event, err := events.New(events.CacophonyBucketUpdate)
		if err != nil {
			return false, err
		}
		event.BucketUpdate = &events.BucketUpdate{
			Tag:       t.TagSuffix,
			GuildID:   env.GuildID,
			Type:      t.Type,
			KeySuffix: key,
			EnvDatas:  values,
		}
		event.GuildID = env.GuildID

		err, recoverable = env.Event.Publisher().Publish(
			context.TODO(),
			event,
		)
		if err != nil {
			if !recoverable {
				event.Logger().Fatal(
					"received unrecoverable error while publishing custom commands alias message",
					zap.Error(err),
				)
			}
			return false, err
		}
	}

	return false, nil
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

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz")

func randomLetters(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func bucketContent(env *models.Env) ([]byte, error) {
	data, err := env.Marshal()
	if err != nil {
		return nil, err
	}

	return append([]byte(randomLetters(10)+"|"), data...), nil
}
