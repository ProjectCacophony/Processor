package filters

import (
	"errors"
	"strconv"

	"gitlab.com/Cacophony/Processor/plugins/automod/interfaces"
	"gitlab.com/Cacophony/Processor/plugins/automod/models"
	"gitlab.com/Cacophony/go-kit/events"
)

type BucketAmount struct{}

func (f BucketAmount) Name() string {
	return "if_bucket_amount"
}

func (f BucketAmount) Args() int {
	return 2
}

func (f BucketAmount) Deprecated() bool {
	return false
}

func (f BucketAmount) NewItem(env *models.Env, args []string) (interfaces.FilterItemInterface, error) {
	if len(args) < 1 {
		return nil, errors.New("too few arguments")
	}

	comparisonType, err := extractAmountComparisonType(args[0])
	if err != nil {
		return nil, err
	}

	amount, err := strconv.Atoi(args[1])
	if err != nil {
		return nil, err
	}

	if amount < 0 {
		return nil, errors.New("amount has to be 0 or greater")
	}

	return &BucketAmountItem{
		Amount:     amount,
		Comparison: comparisonType,
	}, nil
}

func (f BucketAmount) Description() string {
	return "automod.filters.if_bucket_amount"
}

type BucketAmountItem struct {
	Amount     int
	Comparison AmountComparisonType
}

func (f *BucketAmountItem) Match(env *models.Env) bool {
	if env.Event == nil {
		return false
	}

	if env.Event.Type != events.CacophonyBucketUpdate {
		return false
	}

	switch f.Comparison {
	case AmountComparisonLT:
		return len(env.Event.BucketUpdate.EnvDatas) < f.Amount
	case AmountComparisonEQ:
		return len(env.Event.BucketUpdate.EnvDatas) == f.Amount
	case AmountComparisonGT:
		return len(env.Event.BucketUpdate.EnvDatas) > f.Amount
	}

	return false
}
