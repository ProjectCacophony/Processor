package filters

import (
	"errors"
	"strconv"

	"gitlab.com/Cacophony/Processor/plugins/automod/interfaces"
	"gitlab.com/Cacophony/Processor/plugins/automod/models"
	"gitlab.com/Cacophony/go-kit/events"
)

type MentionsCount struct{}

func (f MentionsCount) Name() string {
	return "if_mentions_count"
}

func (f MentionsCount) Args() int {
	return 2
}

func (f MentionsCount) Deprecated() bool {
	return false
}

func (f MentionsCount) NewItem(env *models.Env, args []string) (interfaces.FilterItemInterface, error) {
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

	return &MentionsCountItem{
		Amount:     amount,
		Comparison: comparisonType,
	}, nil
}

func (f MentionsCount) Description() string {
	return "automod.filters.if_mentions_count"
}

type MentionsCountItem struct {
	Amount     int
	Comparison AmountComparisonType
}

func (f *MentionsCountItem) Match(env *models.Env) bool {
	if env.Event == nil {
		return false
	}

	if env.Event.Type != events.MessageCreateType {
		return false
	}

	switch f.Comparison {
	case AmountComparisonLT:
		return len(env.Event.MessageCreate.Mentions) < f.Amount
	case AmountComparisonEQ:
		return len(env.Event.MessageCreate.Mentions) == f.Amount
	case AmountComparisonGT:
		return len(env.Event.MessageCreate.Mentions) > f.Amount
	}

	return false
}
