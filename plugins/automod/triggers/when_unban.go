package triggers

import (
	"gitlab.com/Cacophony/Processor/plugins/automod/interfaces"
	"gitlab.com/Cacophony/Processor/plugins/automod/models"
	"gitlab.com/Cacophony/go-kit/events"
)

type Unban struct{}

func (t Unban) Name() string {
	return "when_unban"
}

func (t Unban) Args() int {
	return 0
}

func (t Unban) Deprecated() bool {
	return false
}

func (t Unban) NewItem(env *models.Env, args []string) (interfaces.TriggerItemInterface, error) {
	return &UnbanItem{}, nil
}

func (t Unban) Description() string {
	return "automod.triggers.when_unban"
}

type UnbanItem struct{}

func (t *UnbanItem) Match(env *models.Env) bool {
	if env.Event == nil {
		return false
	}

	if env.Event.Type != events.GuildBanRemoveType {
		return false
	}

	if env.Event.GuildBanRemove.User.ID == env.Event.BotUserID {
		return false
	}

	env.GuildID = env.Event.GuildBanRemove.GuildID
	env.UserID = append(env.UserID, env.Event.GuildBanRemove.User.ID)

	return true
}
