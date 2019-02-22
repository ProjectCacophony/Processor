package filters

import (
	"errors"
	"strconv"

	"gitlab.com/Cacophony/Processor/plugins/automod/interfaces"
	"gitlab.com/Cacophony/Processor/plugins/automod/models"
	"gitlab.com/Cacophony/go-kit/events"
)

type BucketGT struct {
}

func (f BucketGT) Name() string {
	return "if_bucket_gt"
}

func (f BucketGT) Args() int {
	return 1
}

func (f BucketGT) NewItem(env *models.Env, args []string) (interfaces.FilterItemInterface, error) {
	if len(args) < 1 {
		return nil, errors.New("too few arguments")
	}

	amount, err := strconv.Atoi(args[0])
	if err != nil {
		return nil, err
	}

	return &BucketBTItem{
		Amount: amount,
	}, nil
}

func (f BucketGT) Description() string {
	return "automod.filters.if_bucket_gt"
}

type BucketBTItem struct {
	Amount int
}

func (f *BucketBTItem) Match(env *models.Env) bool {
	if env.Event == nil {
		return false
	}

	if env.Event.Type != events.CacophonyBucketUpdate {
		return false
	}

	return len(env.Event.BucketUpdate.Values) > f.Amount
}
