package handler

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"

	"gitlab.com/Cacophony/go-kit/discord"

	"gitlab.com/Cacophony/Processor/plugins/automod/models"

	"gitlab.com/Cacophony/go-kit/config"
)

const automodLogKey = "cacophony:processor:automod:log-channel-id"

func (h *Handler) getLogChannelIDs() (map[string]string, error) {
	var items []config.Item

	err := h.db.Model(config.Item{}).Where(
		"key = ?",
		automodLogKey,
	).Find(&items).Error
	if err != nil {
		return nil, err
	}

	var valueString string
	list := make(map[string]string)
	for _, item := range items {
		valueString = string(item.Value)
		if valueString == "" {
			continue
		}

		parts := strings.Split(item.Namespace, ":")
		if len(parts) < 2 {
			continue
		}
		list[parts[1]] = valueString
	}

	return list, nil
}

func (h *Handler) postLog(env *models.Env, rule models.Rule) error {
	if env == nil || env.GuildID == "" {
		return nil
	}

	h.logChannelsLock.RLock()
	channelID := h.logChannels[env.GuildID]
	h.logChannelsLock.RUnlock()

	if channelID == "" {
		return nil
	}

	botID, err := env.State.BotForGuild(env.GuildID)
	if err != nil {
		return err
	}

	session, err := discord.NewSession(h.tokens, botID)
	if err != nil {
		return err
	}

	_, err = discord.SendComplexWithVars(
		nil,
		session,
		nil,
		channelID,
		&discordgo.MessageSend{
			Content: fmt.Sprintf("Rule `%s` was triggered.", rule.Name),
		},
		false,
	)
	return err
}
