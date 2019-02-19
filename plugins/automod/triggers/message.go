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

func (t Message) NewItem() interfaces.TriggerItemInterface {
	return &MessageItem{}
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
