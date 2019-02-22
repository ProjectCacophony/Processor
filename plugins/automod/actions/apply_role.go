package actions

import (
	"errors"

	"gitlab.com/Cacophony/Processor/plugins/automod/interfaces"
	"gitlab.com/Cacophony/Processor/plugins/automod/models"
)

type ApplyRole struct {
}

func (t ApplyRole) Name() string {
	return "apply_role"
}

func (t ApplyRole) Args() int {
	return 1
}

func (t ApplyRole) NewItem(env *models.Env, args []string) (interfaces.ActionItemInterface, error) {
	if len(args) < 1 {
		return nil, errors.New("too few arguments")
	}

	role, err := env.State.RoleFromMention(env.GuildID, args[0])
	if err != nil {
		return nil, err
	}

	return &ApplyRoleItem{
		RoleID: role.ID,
	}, nil
}

func (t ApplyRole) Description() string {
	return "automod.actions.apply_role"
}

type ApplyRoleItem struct {
	RoleID string
}

func (t *ApplyRoleItem) Do(env *models.Env) {
	env.Event.Discord().GuildMemberRoleAdd(env.GuildID, env.UserID, t.RoleID) // nolint: errcheck
}
