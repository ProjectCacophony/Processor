package actions

import (
	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/discord"

	"gitlab.com/Cacophony/Processor/plugins/automod/interfaces"
	"gitlab.com/Cacophony/Processor/plugins/automod/models"
)

type BanUser struct {
}

func (t BanUser) Name() string {
	return "ban_user"
}

func (t BanUser) Args() int {
	return 1
}

func (t BanUser) Deprecated() bool {
	return false
}

func (t BanUser) NewItem(env *models.Env, args []string) (interfaces.ActionItemInterface, error) {
	return &BanUserItem{}, nil
}

func (t BanUser) Description() string {
	return "automod.actions.ban_user"
}

type BanUserItem struct {
}

func (t *BanUserItem) Do(env *models.Env) (bool, error) {
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
			discordgo.PermissionBanMembers,
		)
		if err != nil {
			continue
		}

		session, err := discord.NewSession(env.Tokens, botID)
		if err != nil {
			continue
		}

		// TODO: improve Reason
		err = session.Client.GuildBanCreateWithReason(
			env.GuildID, userID, automodReason(env, "Banned"), 0,
		)
		if err != nil {
			return false, err
		}

		doneUserIDs[userID] = true
	}

	return false, nil
}
