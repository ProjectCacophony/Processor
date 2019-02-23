package filters

import (
	"errors"
	"time"

	"gitlab.com/Cacophony/Processor/plugins/automod/interfaces"
	"gitlab.com/Cacophony/Processor/plugins/automod/models"
	"gitlab.com/Cacophony/go-kit/discord"
)

type AccountAge struct {
}

func (f AccountAge) Name() string {
	return "if_account_age"
}

func (f AccountAge) Args() int {
	return 2
}

func (f AccountAge) NewItem(env *models.Env, args []string) (interfaces.FilterItemInterface, error) {
	if len(args) < 1 {
		return nil, errors.New("too few arguments")
	}

	comparisonType, err := extractAmountComparisonType(args[0])
	if err != nil || comparisonType == AmountComparisonEQ {
		return nil, errors.New("invalid Comparison type, try \"<\", or \">\"")
	}

	age, err := time.ParseDuration(args[1])
	if err != nil {
		return nil, err
	}

	return &AccountAgeItem{
		Age:        age,
		Comparison: comparisonType,
	}, nil
}

func (f AccountAge) Description() string {
	return "automod.filters.if_account_age"
}

type AccountAgeItem struct {
	Age        time.Duration
	Comparison AmountComparisonType
}

func (f *AccountAgeItem) Match(env *models.Env) bool {
	var didNotMatch bool

	for _, userID := range env.UserID {
		user, err := env.State.User(userID)
		if err != nil {
			return false
		}

		userTime, err := discord.TimeFromID(user.ID)
		if err != nil {
			return false
		}

		switch f.Comparison {
		case AmountComparisonLT:
			if time.Now().UTC().Add(-f.Age).Before(userTime.UTC()) {
				continue
			}
		case AmountComparisonGT:
			if time.Now().UTC().Add(-f.Age).After(userTime.UTC()) {
				continue
			}
		}

		didNotMatch = true
	}

	return !didNotMatch
}
