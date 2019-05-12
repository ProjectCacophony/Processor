package actions

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/discord"
	"gitlab.com/Cacophony/go-kit/permissions"

	"gitlab.com/Cacophony/Processor/plugins/automod/interfaces"
	"gitlab.com/Cacophony/Processor/plugins/automod/models"
)

type SendMessageTo struct {
}

func (t SendMessageTo) Name() string {
	return "send_message_to"
}

func (t SendMessageTo) Args() int {
	return 2
}

func (t SendMessageTo) NewItem(env *models.Env, args []string) (interfaces.ActionItemInterface, error) {
	if len(args) < 2 {
		return nil, errors.New("too few arguments")
	}

	channel, err := env.State.ChannelFromMention(env.GuildID, args[1])
	if err != nil {
		return nil, err
	}

	return &SendMessageToItem{
		Message:   args[0],
		ChannelID: channel.ID,
	}, nil
}

func (t SendMessageTo) Description() string {
	return "automod.actions.send_message_to"
}

type SendMessageToItem struct {
	Message   string
	ChannelID string
}

func (t *SendMessageToItem) Do(env *models.Env) {
	_, err := env.State.Channel(t.ChannelID)
	if err != nil {
		return
	}

	botID, err := env.State.BotForChannel(
		t.ChannelID,
		permissions.DiscordSendMessages,
	)
	if err != nil {
		return
	}

	session, err := discord.NewSession(env.Tokens, botID)
	if err != nil {
		return
	}

	discord.SendComplexWithVars(
		session,
		nil,
		t.ChannelID,
		&discordgo.MessageSend{
			Content: t.Message,
		},
	)
}
