package triggers

import (
	"gitlab.com/Cacophony/Processor/plugins/automod/interfaces"
	"gitlab.com/Cacophony/Processor/plugins/automod/models"
	"gitlab.com/Cacophony/go-kit/events"
)

type Message struct {
}

func (t Message) Name() string {
	return "when_message"
}

func (t Message) NewItem(env *models.Env) interfaces.TriggerItemInterface {
	return &MessageItem{}
}

func (t Message) Description() string {
	return "automod.triggers.when_message"
}

type MessageItem struct {
}

func (t *MessageItem) Match(env *models.Env) bool {
	if env.Event == nil {
		return false
	}

	if env.Event.Type != events.MessageCreateType {
		return false
	}

	return true
}
