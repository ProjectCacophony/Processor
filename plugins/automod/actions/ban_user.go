package actions

import (
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

func (t BanUser) NewItem(env *models.Env, args []string) (interfaces.ActionItemInterface, error) {
	return &BanUserItem{}, nil
}

func (t BanUser) Description() string {
	return "automod.actions.ban_user"
}

type BanUserItem struct {
}

func (t *BanUserItem) Do(env *models.Env) {
	doneUserIDs := make(map[string]interface{})

	for _, userID := range env.UserID {
		if doneUserIDs[userID] != nil {
			continue
		}

		_, err := env.State.Member(env.GuildID, userID)
		if err != nil {
			continue
		}

		botID, err := env.State.BotForGuild(env.GuildID)
		if err != nil {
			continue
		}

		session, err := discord.NewSession(env.Tokens, botID)
		if err != nil {
			continue
		}

		// nolint: errcheck
		session.GuildBanCreateWithReason(
			env.GuildID, userID, "Banned by Cacophony Automod", 0,
		)
		// TODO: improve Reason

		doneUserIDs[userID] = true
	}
}
