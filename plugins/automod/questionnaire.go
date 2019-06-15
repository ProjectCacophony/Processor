package automod

import (
	"encoding/json"
	"errors"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/Processor/plugins/automod/models"
	"gitlab.com/Cacophony/go-kit/discord"
	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleConfirmUpdateRuleQuestionnaire(event *events.Event) {
	messageID, ok := event.QuestionnaireMatch.Payload["messageID"].(string)
	if !ok || messageID == "" {
		event.Except(errors.New("invalid payload, messageID is empty"))
		return
	}

	ruleData, ok := event.QuestionnaireMatch.Payload["rule"].(string)
	if !ok || len(ruleData) <= 0 {
		event.Except(errors.New("invalid payload, rule data is empty"))
		return
	}

	var rule models.Rule
	err := json.Unmarshal([]byte(ruleData), &rule)
	if err != nil {
		event.Except(err)
		return
	}
	if rule.Name == "" || rule.GuildID == "" {
		event.Except(errors.New("invalid payload, rule name is empty"))
		return
	}

	if event.MessageReactionAdd.MessageID != messageID {
		err = event.Questionnaire().Redo(event)
		event.Except(err)
		return
	}

	switch event.MessageReactionAdd.Emoji.APIName() {
	case "✅":
		err = p.db.Delete(
			models.Rule{}, "guild_id = ? AND name = ?", rule.GuildID, rule.Name,
		).Error
		if err != nil {
			event.Except(err)
			return
		}

		err = p.db.Save(&rule).Error
		if err != nil {
			event.Except(err)
			return
		}

		content := "automote.add.success"

		discord.EditComplexWithVars(
			event.Redis(), event.Discord(), event.Localizations(), &discordgo.MessageEdit{
				Content: &content,
				ID:      messageID,
				Channel: event.MessageReactionAdd.ChannelID,
			}, false,
		)

		err = discord.RemoveReact(
			event.Redis(), event.Discord(),
			event.MessageReactionAdd.ChannelID, messageID,
			event.BotUserID, false,
			"✅",
		)
		if err != nil {
			return
		}
		err = discord.RemoveReact(
			event.Redis(), event.Discord(),
			event.MessageReactionAdd.ChannelID, messageID,
			event.BotUserID, false,
			"❌",
		)
		if err != nil {
			return
		}
		discord.RemoveReact(
			event.Redis(), event.Discord(),
			event.MessageReactionAdd.ChannelID, messageID,
			event.MessageReactionAdd.UserID, false,
			event.MessageReactionAdd.Emoji.APIName(),
		)

		return

	case "❌":

		discord.Delete(
			event.Redis(), event.Discord(), event.MessageReactionAdd.ChannelID, messageID, false,
		)
		return

	}

	err = event.Questionnaire().Redo(event)
	event.Except(err)
}
