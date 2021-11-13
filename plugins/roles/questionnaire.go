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

		err = p.db.Where("category_id = ?", category.ID).Delete(&Role{}).Error
		if err != nil {
			event.Except(err)
			return
		}

		err = p.db.Delete(category).Error
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

func (p *Plugin) handleConfirmCreateRole(event *events.Event) {
	switch event.MessageReactionAdd.Emoji.APIName() {
	case "✅":
		messageID, ok := event.QuestionnaireMatch.Payload["messageID"].(string)
		if !ok || messageID == "" {
			event.Except(errors.New("invalid payload, messageID is empty"))
			return
		}

		roleName, ok := event.QuestionnaireMatch.Payload["roleName"].(string)
		if !ok || roleName == "" {
			event.Except(errors.New("invalid payload, roleName is empty"))
			return
		}

		roleData, ok := event.QuestionnaireMatch.Payload["role"].(string)
		if !ok {
			event.Except(errors.New("invalid payload, role is empty"))
			return
		}

		var role Role
		err := json.Unmarshal([]byte(roleData), &role)
		if err != nil {
			event.Except(err)
			return
		}

		newServerrole, err := event.Discord().Client.GuildRoleCreate(event.GuildID)
		if err != nil {
			event.Except(err)
			return
		}

		newServerrole, err = event.Discord().Client.GuildRoleEdit(event.GuildID, newServerrole.ID, roleName, 0, false, 0, false)
		if err != nil {
			event.Except(err)
			return
		}

		role.ServerRoleID = newServerrole.ID

		err = p.db.Save(&role).Error
		if err != nil {
			event.Except(err)
			return
		}

		discord.Delete(
			event.Redis(), event.Discord(), event.MessageReactionAdd.ChannelID, messageID, false,
		)

		_, err = event.Send(event.ChannelID, "roles.role.creation",
			"roleName", newServerrole.Name,
		)
		if err != nil {
			event.Except(err)
			return
		}

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
