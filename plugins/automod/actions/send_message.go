package actions

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/discord"
	"gitlab.com/Cacophony/go-kit/permissions"

	"gitlab.com/Cacophony/Processor/plugins/automod/interfaces"
	"gitlab.com/Cacophony/Processor/plugins/automod/models"
)

type SendMessage struct {
}

func (t SendMessage) Name() string {
	return "send_message"
}

func (t SendMessage) Args() int {
	return 1
}

func (t SendMessage) Deprecated() bool {
	return false
}

func (t SendMessage) NewItem(env *models.Env, args []string) (interfaces.ActionItemInterface, error) {
	if len(args) < 1 {
		return nil, errors.New("too few arguments")
	}

	return &SendMessageItem{
		Message: args[0],
	}, nil
}

func (t SendMessage) Description() string {
	return "automod.actions.send_message"
}

type SendMessageItem struct {
	Message string
}

func (t *SendMessageItem) Do(env *models.Env) (bool, error) {
	doneChannelIDs := make(map[string]interface{})

	var messages []*discordgo.Message
	for _, channelID := range env.ChannelID {
		if doneChannelIDs[channelID] != nil {
			continue
		}

		_, err := env.State.Channel(channelID)
		if err != nil {
			continue
		}

		botID, err := env.State.BotForChannel(
			channelID,
			permissions.DiscordSendMessages,
		)
		if err != nil {
			continue
		}

		session, err := discord.NewSession(env.Tokens, botID)
		if err != nil {
			continue
		}

		messages, err = discord.SendComplexWithVars(
			session,
			nil,
			channelID,
			discord.MessageCodeToMessage(ReplaceText(env, t.Message)),
		)
		if err != nil {
			return false, err
		}
		for _, message := range messages {
			env.Messages = append(env.Messages, models.NewEnvMessage(message))
		}

		doneChannelIDs[channelID] = true
	}

	return false, nil
}
