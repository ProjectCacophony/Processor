package actions

import (
	"errors"

	"gitlab.com/Cacophony/Processor/plugins/automod/interfaces"
	"gitlab.com/Cacophony/Processor/plugins/automod/models"
)

type SendMessage struct {
}

func (t SendMessage) Name() string {
	return "send_message"
}

func (t SendMessage) Args() int {
	return 1
}

func (t SendMessage) NewItem(env *models.Env, args []string) (interfaces.ActionItemInterface, error) {
	if len(args) < 1 {
		return nil, errors.New("too few arguments")
	}

	return &SendMessageItem{
		Message: args[0],
	}, nil
}

func (t SendMessage) Description() string {
	return "automod.actions.send_message"
}

type SendMessageItem struct {
	Message string
}

func (t *SendMessageItem) Do(env *models.Env) {
	env.Event.Respond(t.Message) // nolint: errcheck
}
