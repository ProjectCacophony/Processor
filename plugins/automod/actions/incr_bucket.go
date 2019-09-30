package actions

import (
	"errors"
	"math/rand"
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

func (t *IncrBucketItem) Do(env *models.Env) error {
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

	var values []string

	for _, key := range keys {

		for i := 0; i < t.Amount; i++ {

			valueList, _ := bucket.AddWithValue(
				env.Redis,
				key,
				bucketContent(env),
				t.Decay,
			)

			values = make([]string, len(valueList))

			for i, value := range valueList {
				stringValue, ok := value.Member.(string)
				if !ok {
					continue
				}

				values[i] = stringValue
			}
		}

		// TODO: publish event?
		// TODO: async
		env.Handler.Handle(&events.Event{
			GuildID: env.GuildID,
			Type:    events.CacophonyBucketUpdate,
			BucketUpdate: &events.BucketUpdate{
				Tag:       t.TagSuffix,
				GuildID:   env.GuildID,
				Values:    values,
				Type:      t.Type,
				KeySuffix: key,
			},
		})
	}

	return nil
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

func bucketContent(env *models.Env) string {
	// TODO: marshal bucket content in a sane way

	var userText, channelText, messageText string

	for _, user := range env.UserID {
		if strings.Contains(userText, user+";") {
			continue
		}

		userText += user + ";"
	}

	for _, channel := range env.ChannelID {
		if strings.Contains(channelText, channel+";") {
			continue
		}

		channelText += channel + ";"
	}

	for _, message := range env.Messages {
		if strings.Contains(messageText, message.ID+":"+message.ChanneID+";") {
			continue
		}

		messageText += message.ID + ":" + message.ChanneID + ";"
	}

	return env.GuildID + "|" +
		channelText + "|" +
		userText + "|" +
		messageText + "|" +
		randomLetters(10)
}
