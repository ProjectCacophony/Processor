package actions

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/discord"

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

func (t ApplyRole) Deprecated() bool {
	return false
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

func (t *ApplyRoleItem) Do(env *models.Env) (bool, error) {
	doneUserIDs := make(map[string]interface{})

	for _, userID := range env.UserID {
		if doneUserIDs[userID] != nil {
			continue
		}

		_, err := env.State.Member(env.GuildID, userID)
		if err != nil {
			continue
		}

		botID, err := env.State.BotForGuild(
			env.GuildID,
			discordgo.PermissionManageRoles,
		)
		if err != nil {
			continue
		}

		session, err := discord.NewSession(env.Tokens, botID)
		if err != nil {
			continue
		}

		err = session.Client.GuildMemberRoleAdd(env.GuildID, userID, t.RoleID)
		if err != nil {
			return false, err
		}

		doneUserIDs[userID] = true
	}

	return false, nil
}
