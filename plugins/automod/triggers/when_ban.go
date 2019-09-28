package triggers

import (
	"gitlab.com/Cacophony/Processor/plugins/automod/interfaces"
	"gitlab.com/Cacophony/Processor/plugins/automod/models"
	"gitlab.com/Cacophony/go-kit/events"
)

type Ban struct {
}

func (t Ban) Name() string {
	return "when_ban"
}

func (t Ban) Args() int {
	return 0
}

func (t Ban) NewItem(env *models.Env, args []string) (interfaces.TriggerItemInterface, error) {
	return &BanItem{}, nil
}

func (t Ban) Description() string {
	return "automod.triggers.when_ban"
}

type BanItem struct {
}

func (t *BanItem) Match(env *models.Env) bool {
	if env.Event == nil {
		return false
	}

	if env.Event.Type != events.GuildBanAddType {
		return false
	}

	if env.Event.GuildBanAdd.User.ID == env.Event.BotUserID {
		return false
	}

	env.GuildID = env.Event.GuildBanAdd.GuildID
	env.UserID = append(env.UserID, env.Event.GuildBanAdd.User.ID)

	return true
}
