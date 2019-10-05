package actions

import (
	"context"
	"errors"
	"time"

	"gitlab.com/Cacophony/go-kit/events"

	"gitlab.com/Cacophony/Processor/plugins/automod/interfaces"
	"gitlab.com/Cacophony/Processor/plugins/automod/models"
)

type Wait struct {
}

func (t Wait) Name() string {
	return "wait"
}

func (t Wait) Args() int {
	return 1
}

func (t Wait) Deprecated() bool {
	return false
}

func (t Wait) NewItem(env *models.Env, args []string) (interfaces.ActionItemInterface, error) {
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

func (t Wait) Description() string {
	return "automod.actions.wait"
}

type WaitItem struct {
	Duration time.Duration
}

func (t *WaitItem) Do(env *models.Env) (bool, error) {
	var newActions []models.RuleAction
	for i, action := range env.Rule.Actions {
		if action.Name != (Wait{}).Name() {
			continue
		}

		newActions = env.Rule.Actions[i+1:]
		break
	}
	env.Rule.Filters = nil
	env.Rule.Actions = newActions

	envData, err := env.Marshal()
	if err != nil {
		return true, err
	}

	event, err := events.New(events.CacophonyAutomodWait)
	if err != nil {
		return true, err
	}

	event.AutomodWait = &events.AutomodWait{
		EnvData: envData,
	}
	event.GuildID = env.GuildID

	err = env.Event.Publisher().PublishAt(
		context.TODO(),
		event,
		time.Now().Add(t.Duration),
	)
	if err != nil {
		return true, err
	}

	return true, nil
}
