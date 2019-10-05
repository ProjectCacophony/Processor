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

func (t Message) Args() int {
	return 0
}

func (t Message) Deprecated() bool {
	return false
}

func (t Message) NewItem(env *models.Env, args []string) (interfaces.TriggerItemInterface, error) {
	return &MessageItem{}, nil
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

	if env.Event.MessageCreate.Author.Bot {
		return false
	}

	env.GuildID = env.Event.MessageCreate.GuildID
	env.UserID = append(env.UserID, env.Event.MessageCreate.Author.ID)
	env.ChannelID = append(env.ChannelID, env.Event.MessageCreate.ChannelID)
	env.Messages = append(env.Messages, models.NewEnvMessage(env.Event.MessageCreate.Message))

	return true
}
