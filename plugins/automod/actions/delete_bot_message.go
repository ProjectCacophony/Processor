package actions

import (
	"gitlab.com/Cacophony/Processor/plugins/automod/interfaces"
	"gitlab.com/Cacophony/Processor/plugins/automod/models"
	"gitlab.com/Cacophony/go-kit/discord"
	"gitlab.com/Cacophony/go-kit/permissions"
)

type DeleteBotMessage struct{}

func (t DeleteBotMessage) Name() string {
	return "delete_bot_message"
}

func (t DeleteBotMessage) Args() int {
	return 0
}

func (t DeleteBotMessage) Deprecated() bool {
	return false
}

func (t DeleteBotMessage) NewItem(env *models.Env, args []string) (interfaces.ActionItemInterface, error) {
	return &DeleteBotMessageItem{}, nil
}

func (t DeleteBotMessage) Description() string {
	return "automod.actions.delete_bot_message"
}

type DeleteBotMessageItem struct{}

func (t *DeleteBotMessageItem) Do(env *models.Env) (bool, error) {
	// TODO: group messages by channel ID, use bulk delete endpoint

	doneMessageIDs := make(map[string]interface{})

	for _, message := range env.Messages {
		if message == nil {
			continue
		}
		if !message.Bot {
			continue
		}

		if doneMessageIDs[message.ID] != nil {
			continue
		}

		botID, err := env.State.BotForChannel(
			message.ChannelID,
			permissions.DiscordManageMessages,
		)
		if err != nil {
			continue
		}

		session, err := discord.NewSession(env.Tokens, botID)
		if err != nil {
			continue
		}

		err = discord.Delete(nil, session, message.ChannelID, message.ID, false)
		if err != nil {
			return false, err
		}

		doneMessageIDs[message.ID] = true
	}

	return false, nil
}
