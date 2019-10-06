package roles

import (
	"encoding/json"
	"errors"

	"gitlab.com/Cacophony/go-kit/discord"
	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleConfirmCategoryDelete(event *events.Event) {

	switch event.MessageReactionAdd.Emoji.APIName() {
	case "✅":
		messageID, ok := event.QuestionnaireMatch.Payload["messageID"].(string)
		if !ok || messageID == "" {
			event.Except(errors.New("invalid payload, messageID is empty"))
			return
		}

		categoryData, ok := event.QuestionnaireMatch.Payload["category"].(string)
		if !ok || len(categoryData) <= 0 {
			event.Except(errors.New("invalid payload, category data is empty"))
			return
		}

		var category Category
		err := json.Unmarshal([]byte(categoryData), &category)
		if err != nil {
			event.Except(err)
			return
		}
		if category.Name == "" || category.GuildID == "" {
			event.Except(errors.New("invalid payload, category name is empty"))
			return
		}

		err = p.db.Delete(category.Roles).Delete(category).Error
		if err != nil {
			event.Except(err)
			return
		}

		event.Send(event.ChannelID, "roles.category.deleted",
			"category", category,
		)

		discord.Delete(
			event.Redis(), event.Discord(), event.MessageReactionAdd.ChannelID, messageID, false,
		)

		return

	case "❌":
		messageID, ok := event.QuestionnaireMatch.Payload["messageID"].(string)
		if !ok || messageID == "" {
			event.Except(errors.New("invalid payload, messageID is empty"))
			return
		}

		discord.Delete(
			event.Redis(), event.Discord(), event.MessageReactionAdd.ChannelID, messageID, false,
		)
		return

	}

	err := event.Questionnaire().Redo(event)
	event.Except(err)
}
