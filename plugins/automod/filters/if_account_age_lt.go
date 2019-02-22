package filters

import (
	"errors"
	"time"

	"gitlab.com/Cacophony/go-kit/discord"

	"gitlab.com/Cacophony/Processor/plugins/automod/interfaces"
	"gitlab.com/Cacophony/Processor/plugins/automod/models"
)

type AccountAgeLT struct {
}

func (f AccountAgeLT) Name() string {
	return "if_account_age_lt"
}

func (f AccountAgeLT) Args() int {
	return 1
}

func (f AccountAgeLT) NewItem(env *models.Env, args []string) (interfaces.FilterItemInterface, error) {
	if len(args) < 1 {
		return nil, errors.New("too few arguments")
	}

	age, err := time.ParseDuration(args[0])
	if err != nil {
		return nil, err
	}

	return &AccountAgeLTItem{
		age: age,
	}, nil
}

func (f AccountAgeLT) Description() string {
	return "automod.filters.if_account_age_lt"
}

type AccountAgeLTItem struct {
	age time.Duration
}

func (f *AccountAgeLTItem) Match(env *models.Env) bool {
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

		if time.Now().UTC().Add(-f.age).Before(userTime.UTC()) {
			continue
		}

		didNotMatch = true
	}

	return !didNotMatch
}
