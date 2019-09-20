package triggers

import (
	"gitlab.com/Cacophony/Processor/plugins/automod/interfaces"
	"gitlab.com/Cacophony/Processor/plugins/automod/models"
	"gitlab.com/Cacophony/go-kit/events"
)

type Leave struct {
}

func (t Leave) Name() string {
	return "when_leave"
}

func (t Leave) Args() int {
	return 0
}

func (t Leave) NewItem(env *models.Env, args []string) (interfaces.TriggerItemInterface, error) {
	return &LeaveItem{}, nil
}

func (t Leave) Description() string {
	return "automod.triggers.when_leave"
}

type LeaveItem struct {
}

func (t *LeaveItem) Match(env *models.Env) bool {
	if env.Event == nil {
		return false
	}

	if env.Event.Type != events.GuildMemberRemoveType {
		return false
	}

	if env.Event.GuildMemberRemove.User.ID == env.Event.BotUserID {
		return false
	}

	env.GuildID = env.Event.GuildMemberRemove.GuildID
	env.UserID = append(env.UserID, env.Event.GuildMemberRemove.User.ID)

	return true
}
