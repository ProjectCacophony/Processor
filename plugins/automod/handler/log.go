package handler

import (
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
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

func (h *Handler) logRun(env *models.Env, rule models.Rule, runError error) error {
	messageIDs := make([]string, len(env.Messages))
	for i, message := range env.Messages {
		messageIDs[i] = message.ChanneID + ":" + message.ID
	}

	entry := &models.LogEntry{
		GuildID:    env.GuildID,
		RuleID:     rule.ID,
		ChannelIDs: env.ChannelID,
		UserIDs:    env.UserID,
		MessageIDs: messageIDs,
	}
	if runError != nil {
		entry.ErrorMessage = runError.Error()
	}

	err := h.db.Save(&entry).Error
	if err != nil {
		return err
	}

	err = h.db.Model(models.Rule{}).
		Where("id = ?", rule.ID).
		Update("runs", gorm.Expr("runs + 1")).
		Error
	if err != nil {
		return err
	}

	if rule.Silent {
		return nil
	}

	return h.postLog(env, rule, runError)
}

func (h *Handler) postLog(env *models.Env, rule models.Rule, runError error) error {
	if env == nil || env.GuildID == "" {
		return nil
	}

	h.logChannelsLock.RLock()
	channelID := h.logChannels[env.GuildID]
	h.logChannelsLock.RUnlock()

	if channelID == "" {
		return nil
	}

	botID, err := env.State.BotForChannel(
		channelID,
		permissions.DiscordSendMessages,
		permissions.DiscordEmbedLinks,
	)
	if err != nil {
		return err
	}

	session, err := discord.NewSession(h.tokens, botID)
	if err != nil {
		return err
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

	description := "✅"
	if runError != nil {
		description = "⚠ " + runError.Error()
	}
	description += fmt.Sprintf("\nName: `%s`", rule.Name)

	_, err = discord.SendComplexWithVars(
		session,
		nil,
		channelID,
		&discordgo.MessageSend{
			Embed: &discordgo.MessageEmbed{
				Title:       "Automod: Rule Triggered",
				Description: description,
				Timestamp:   time.Now().Format(time.RFC3339),
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
	)
	return err
}
