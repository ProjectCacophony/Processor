package actions

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"gitlab.com/Cacophony/go-kit/bucket"

	"gitlab.com/Cacophony/Processor/plugins/automod/interfaces"
	"gitlab.com/Cacophony/Processor/plugins/automod/models"
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
		Tag:    bucketTag(args[0]),
		Decay:  decay,
		Amount: amount,
	}, nil
}

func (t IncrBucket) Description() string {
	return "automod.actions.incr_bucket"
}

type IncrBucketItem struct {
	Tag    string
	Decay  time.Duration
	Amount int
}

func (t *IncrBucketItem) Do(env *models.Env) {
	// TODO: add support for amount
	amount, _ := bucket.Add(
		env.Redis,
		t.Tag,
		t.Decay,
	)
	fmt.Println("bucket", t.Tag, "is now at", amount)
}

func bucketTag(tag string) string {
	return "automod:" + tag
}
