package actions

import (
	"gitlab.com/Cacophony/Processor/plugins/automod/interfaces"
	"gitlab.com/Cacophony/Processor/plugins/automod/models"
)

type SendMessage struct {
}

func (t SendMessage) Name() string {
	return "send_message"
}

func (t SendMessage) NewItem(env *models.Env, value string) (interfaces.ActionItemInterface, error) {
	return &SendMessageItem{
		Message: value,
	}, nil
}

type SendMessageItem struct {
	Message string
}

func (t *SendMessageItem) Do(env *models.Env) {
	env.Event.Respond(t.Message) // nolint: errcheck
}
