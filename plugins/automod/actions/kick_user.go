package actions

import (
	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/discord"

	"gitlab.com/Cacophony/Processor/plugins/automod/interfaces"
	"gitlab.com/Cacophony/Processor/plugins/automod/models"
)

type KickUser struct {
}

func (t KickUser) Name() string {
	return "kick_user"
}

func (t KickUser) Args() int {
	return 0
}

func (t KickUser) NewItem(env *models.Env, args []string) (interfaces.ActionItemInterface, error) {
	return &KickUserItem{}, nil
}

func (t KickUser) Description() string {
	return "automod.actions.kick_user"
}

type KickUserItem struct {
}

func (t *KickUserItem) Do(env *models.Env) {
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
			discordgo.PermissionKickMembers,
		)
		if err != nil {
			continue
		}

		session, err := discord.NewSession(env.Tokens, botID)
		if err != nil {
			continue
		}

		session.Client.GuildMemberDeleteWithReason(
			env.GuildID, userID, "Kicked by Cacophony Automod",
		)
		// TODO: improve Reason

		doneUserIDs[userID] = true
	}
}
