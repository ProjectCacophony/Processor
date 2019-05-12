package actions

import (
	"errors"
	"regexp"
	"strings"

	"gitlab.com/Cacophony/go-kit/discord"
	"gitlab.com/Cacophony/go-kit/permissions"

	"gitlab.com/Cacophony/Processor/plugins/automod/interfaces"
	"gitlab.com/Cacophony/Processor/plugins/automod/models"
)

var emojiRegex = regexp.MustCompile(`^[\x{00A0}-\x{1F9EF}]|<(a)?:[^<>:]+:[0-9]+>$`)

type React struct {
}

func (t React) Name() string {
	return "react"
}

func (t React) Args() int {
	return 1
}

func (t React) NewItem(env *models.Env, args []string) (interfaces.ActionItemInterface, error) {
	if len(args) < 1 {
		return nil, errors.New("too few arguments")
	}

	if !emojiRegex.MatchString(args[0]) {
		return nil, errors.New("invalid emoji")
	}

	// TODO: confirm that we have access to the emoji

	return &ReactItem{
		Reaction: strings.Trim(args[0], "<>"),
	}, nil
}

func (t React) Description() string {
	return "automod.actions.react"
}

type ReactItem struct {
	Reaction string
}

func (t *ReactItem) Do(env *models.Env) {
	doneMessageIDs := make(map[string]interface{})

	for _, message := range env.Messages {
		if doneMessageIDs[message.ID] != nil {
			continue
		}

		_, err := env.State.Channel(message.ChanneID)
		if err != nil {
			continue
		}

		botID, err := env.State.BotForChannel(
			message.ChanneID,
			permissions.DiscordAddReactions,
		)
		if err != nil {
			continue
		}

		session, err := discord.NewSession(env.Tokens, botID)
		if err != nil {
			continue
		}

		discord.React(
			nil,
			session,
			message.ChanneID,
			message.ID,
			false,
			t.Reaction,
		)

		doneMessageIDs[message.ID] = true
	}
}
