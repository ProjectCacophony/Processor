package handler

import (
	"fmt"
	"strings"
	"time"

	"gitlab.com/Cacophony/go-kit/permissions"

	"github.com/bwmarrin/discordgo"

	"gitlab.com/Cacophony/go-kit/discord"

	"gitlab.com/Cacophony/Processor/plugins/automod/models"

	"gitlab.com/Cacophony/go-kit/config"
)

const AutomodLogKey = "cacophony:processor:automod:log-channel-id"

func (h *Handler) getLogChannelIDs() (map[string]string, error) {
	var items []config.Item

	err := h.db.Model(config.Item{}).Where(
		"key = ?",
		AutomodLogKey,
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

	if !permissions.And(
		permissions.DiscordSendMessages,
		permissions.DiscordEmbedLinks,
	).Match(h.state, botID, channelID, false) {
		return nil
	}

	usersText := strings.Join(env.UserID, ">, <@")
	if usersText != "" {
		usersText = "<@" + usersText + ">"
	}

	channelsText := strings.Join(env.ChannelID, ">, <#")
	if channelsText != "" {
		channelsText = "<#" + channelsText + ">"
	}

	var actionsText string
	for _, action := range rule.Actions {
		actionsText += "`" + action.Name + "`, "
	}
	if rule.Stop {
		actionsText += "stop, "
	}
	if rule.Silent {
		actionsText += "silent, "
	}
	actionsText = strings.TrimRight(actionsText, ", ")

	_, err = discord.SendComplexWithVars(
		nil,
		session,
		nil,
		channelID,
		&discordgo.MessageSend{
			Embed: &discordgo.MessageEmbed{
				URL:         "",
				Type:        "",
				Title:       "Automod: Rule Triggered",
				Description: fmt.Sprintf("Name: `%s`", rule.Name),
				Timestamp:   time.Now().Format(time.RFC3339),
				Color:       0,
				Footer:      nil,
				Image:       nil,
				Thumbnail:   nil,
				Video:       nil,
				Provider:    nil,
				Author:      nil,
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "User(s)",
						Value:  usersText,
						Inline: true,
					},
					{
						Name:   "Channel(s)",
						Value:  channelsText,
						Inline: true,
					},
					{
						Name:   "Action(s)",
						Value:  actionsText,
						Inline: false,
					},
				},
			},
		},
		false,
	)
	return err
}
