package actions

import (
	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/Processor/plugins/automod/interfaces"
	"gitlab.com/Cacophony/Processor/plugins/automod/models"
	"gitlab.com/Cacophony/go-kit/discord"
	"gitlab.com/Cacophony/go-kit/permissions"
)

type Publish struct{}

func (t Publish) Name() string {
	return "publish"
}

func (t Publish) Args() int {
	return 0
}

func (t Publish) Deprecated() bool {
	return false
}

func (t Publish) NewItem(env *models.Env, args []string) (interfaces.ActionItemInterface, error) {
	return &PublishItem{}, nil
}

func (t Publish) Description() string {
	return "automod.actions.publish"
}

type PublishItem struct{}

func (t *PublishItem) Do(env *models.Env) (bool, error) {
	doneMessageIDs := make(map[string]interface{})

	for _, message := range env.Messages {
		if message == nil {
			continue
		}

		if doneMessageIDs[message.ID] != nil {
			continue
		}

		messageChannel, err := env.State.Channel(message.ChannelID)
		if err != nil {
			continue
		}
		if messageChannel.Type != discordgo.ChannelTypeGuildNews {
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

		_, err = session.Client.ChannelMessageCrosspost(message.ChannelID, message.ID)
		if err != nil {
			return false, err
		}

		doneMessageIDs[message.ID] = true
	}

	return false, nil
}
