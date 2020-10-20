package triggers

import (
	"gitlab.com/Cacophony/Processor/plugins/automod/interfaces"
	"gitlab.com/Cacophony/Processor/plugins/automod/models"
	"gitlab.com/Cacophony/go-kit/events"
)

type BotMessage struct {
}

func (t BotMessage) Name() string {
	return "when_bot_message"
}

func (t BotMessage) Args() int {
	return 0
}

func (t BotMessage) Deprecated() bool {
	return false
}

func (t BotMessage) NewItem(env *models.Env, args []string) (interfaces.TriggerItemInterface, error) {
	return &BotMessageItem{}, nil
}

func (t BotMessage) Description() string {
	return "automod.triggers.when_bot_Message"
}

type BotMessageItem struct {
}

var botMessageAllowedActions = map[string]bool{
	"publish": true,
}

func (t *BotMessageItem) Match(env *models.Env) bool {
	if env.Event == nil {
		return false
	}

	if env.Event.Type != events.MessageCreateType {
		return false
	}

	if !env.Event.MessageCreate.Author.Bot {
		return false
	}

	// ignore matches if one action is not a valid action
	if env.Rule == nil || len(env.Rule.Actions) <= 0 {
		return false
	}
	for _, action := range env.Rule.Actions {
		if !botMessageAllowedActions[action.Name] {
			return false
		}
	}

	env.GuildID = env.Event.MessageCreate.GuildID
	env.UserID = append(env.UserID, env.Event.MessageCreate.Author.ID)
	env.ChannelID = append(env.ChannelID, env.Event.MessageCreate.ChannelID)
	env.Messages = append(env.Messages, models.NewEnvMessage(env.Event.MessageCreate.Message))

	return true
}
