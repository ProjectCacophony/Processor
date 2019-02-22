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
	return 3
}

func (t IncrBucket) NewItem(env *models.Env, args []string) (interfaces.ActionItemInterface, error) {
	if len(args) < 3 {
		return nil, errors.New("too few arguments")
	}

	decay, err := time.ParseDuration(args[1])
	if err != nil {
		return nil, err
	}

	amount, err := strconv.Atoi(args[2])
	if err != nil {
		return nil, err
	}
	if amount < 1 {
		return nil, errors.New("amount has to be bigger than 0")
	}

	// TODO: sanitize bucket tag

	return &IncrBucketItem{
		GuildID: env.GuildID,
		Tag:     bucketTag(env.GuildID, args[0]),
		Decay:   decay,
		Amount:  amount,
	}, nil
}

func (t IncrBucket) Description() string {
	return "automod.actions.incr_bucket"
}

type IncrBucketItem struct {
	Tag     string
	Decay   time.Duration
	Amount  int
	GuildID string
}

func (t *IncrBucketItem) Do(env *models.Env) {
	// TODO: add support for amount
	valueList, _ := bucket.AddWithValue(
		env.Redis,
		t.Tag,
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
	env.Handler.Handle(&events.Event{
		Type: events.CacophonyBucketUpdate,
		BucketUpdate: &events.BucketUpdate{
			Tag:     t.Tag,
			GuildID: t.GuildID,
			Values:  values,
		},
	})
}

func bucketTag(guildID, tag string) string {
	return "automod:" + guildID + ":" + tag
}

func bucketContent(env *models.Env) string {
	return env.GuildID + "|" +
		strings.Join(env.ChannelID, ";") + "|" +
		strings.Join(env.UserID, ";") + "|" +
		strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
}

// TODO: allow user channel, or user specific buckets (change bucket tag?)
