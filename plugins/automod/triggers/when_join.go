package triggers

import (
	"gitlab.com/Cacophony/Processor/plugins/automod/interfaces"
	"gitlab.com/Cacophony/Processor/plugins/automod/models"
	"gitlab.com/Cacophony/go-kit/events"
)

type Join struct {
}

func (t Join) Name() string {
	return "when_join"
}

func (t Join) Args() int {
	return 0
}

func (t Join) NewItem(env *models.Env, args []string) (interfaces.TriggerItemInterface, error) {
	return &JoinItem{}, nil
}

func (t Join) Description() string {
	return "automod.triggers.when_join"
}

type JoinItem struct {
}

func (t *JoinItem) Match(env *models.Env) bool {
	if env.Event == nil {
		return false
	}

	if env.Event.Type != events.GuildMemberAddType {
		return false
	}

	if env.Event.GuildMemberAdd.User.Bot {
		return false
	}

	env.GuildID = env.Event.GuildMemberAdd.GuildID
	env.UserID = append(env.UserID, env.Event.GuildMemberAdd.User.ID)

	return true
}
