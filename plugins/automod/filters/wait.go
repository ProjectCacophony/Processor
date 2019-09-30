package filters

import (
	"context"
	"errors"
	"time"

	"gitlab.com/Cacophony/Processor/plugins/automod/interfaces"
	"gitlab.com/Cacophony/Processor/plugins/automod/models"
	"gitlab.com/Cacophony/go-kit/events"
)

type Wait struct {
}

func (f Wait) Name() string {
	return "wait"
}

func (f Wait) Args() int {
	return 1
}

func (f Wait) Deprecated() bool {
	return false
}

func (f Wait) NewItem(env *models.Env, args []string) (interfaces.FilterItemInterface, error) {
	if len(args) < 1 {
		return nil, errors.New("too few arguments")
	}

	duration, err := time.ParseDuration(args[0])
	if err != nil {
		return nil, err
	}

	if duration < 1*time.Second {
		return nil, errors.New("duration has to be equal to or greater than 1 second")
	}
	if duration > 24*time.Hour {
		return nil, errors.New("duration has to be equal to or less than 1 day")
	}

	return &WaitItem{
		Duration: duration,
	}, nil
}

func (f Wait) Description() string {
	return "automod.filters.wait"
}

type WaitItem struct {
	Duration time.Duration
}

func (f *WaitItem) Match(env *models.Env) bool {
	event, err := events.New(events.CacophonyAutomodWait)
	if err != nil {
		env.Event.ExceptSilent(err)
		return false
	}

	var newFilters []models.RuleFilter
	for i, filter := range env.Rule.Filters {
		if filter.Name != (Wait{}).Name() {
			continue
		}

		newFilters = env.Rule.Filters[i+1:]
		break
	}
	env.Rule.Filters = newFilters

	envData, err := env.Marshal()
	if err != nil {
		env.Event.ExceptSilent(err)
		return false
	}

	event.AutomodWait = &events.AutomodWait{
		Payload: envData,
	}
	event.GuildID = env.GuildID

	err = env.Event.Publisher().PublishAt(
		context.TODO(),
		event,
		time.Now().Add(f.Duration),
	)
	if err != nil {
		env.Event.ExceptSilent(err)
		return false
	}

	// to stop further execution we say the event did not match
	return false
}
