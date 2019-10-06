package roles

import (
	"encoding/json"
	"strconv"

	"gitlab.com/Cacophony/go-kit/discord"
	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) createCategory(event *events.Event) {
	if len(event.Fields()) < 4 {
		event.Respond("common.invalid-params")
		return
	}

	name := event.Fields()[3]

	if name == "" {
		event.Respond("roles.category.no-name")
		return
	}

	existingCategory, err := p.getCategoryByName(name, event.GuildID)
	if err != nil {
		event.Except(err)
		return
	}
	if existingCategory.Name != "" {
		event.Respond("roles.category.exists")
		return
	}

	var channelID string
	if len(event.Fields()) >= 5 {
		channel, err := event.State().ChannelFromMention(event.GuildID, event.Fields()[4])
		if err != nil {
			event.Except(err)
			return
		}
		channelID = channel.ID
	}

	limit := 0
	if len(event.Fields()) >= 6 {
		limit, err = strconv.Atoi(event.Fields()[5])
		if err != nil {
			event.Respond("roles.category.limit-not-number")
			return
		}
	}

	category := &Category{
		GuildID:   event.GuildID,
		ChannelID: channelID,
		Name:      name,
		Limit:     limit,
		Enabled:   true,
	}

	err = p.db.Save(category).Error
	if err != nil {
		event.Except(err)
		return
	}

	event.Respond("roles.category.creation",
		"category", category,
	)
}

func (p *Plugin) updateCategory(event *events.Event) {
	if len(event.Fields()) < 5 {
		event.Respond("common.invalid-params")
		return
	}

	currentName := event.Fields()[3]
	existingCategory, err := p.getCategoryByName(currentName, event.GuildID)
	if err != nil {
		event.Except(err)
		return
	}
	if existingCategory.Name == "" {
		event.Respond("roles.category.does-not-exist")
		return
	}

	name := event.Fields()[4]

	if name == "" {
		event.Respond("roles.category.no-name")
		return
	}

	var channelID string
	if len(event.Fields()) >= 6 {
		channel, err := event.State().ChannelFromMention(event.GuildID, event.Fields()[5])
		if err != nil {
			event.Except(err)
			return
		}
		channelID = channel.ID
	}

	limit := 0
	if len(event.Fields()) >= 7 {
		limit, err = strconv.Atoi(event.Fields()[6])
		if err != nil {
			event.Respond("roles.category.limit-not-number")
			return
		}
	}

	existingCategory.ChannelID = channelID
	existingCategory.Name = name
	existingCategory.Limit = limit

	err = p.db.Save(existingCategory).Error
	if err != nil {
		event.Except(err)
		return
	}

	event.Respond("roles.category.update",
		"category", existingCategory,
	)
}

func (p *Plugin) deleteCategory(event *events.Event) {
	if len(event.Fields()) < 4 {
		event.Respond("common.invalid-params")
		return
	}

	category, err := p.getCategoryByName(event.Fields()[3], event.GuildID)
	if err != nil {
		event.Except(err)
		return
	}

	if category.Name == "" {
		event.Respond("roles.category.does-not-exist")
		return
	}

	if len(category.Roles) > 0 {
		messages, err := event.Respond("roles.category.confirm-delete", "rolesCount", len(category.Roles))
		if err != nil {
			event.Except(err)
			return
		}

		categoryData, err := json.Marshal(&category)
		if err != nil {
			event.Except(err)
			return
		}

		err = event.Questionnaire().Register(
			confirmCategoryDeleteKey,
			events.QuestionnaireFilter{
				GuildID:   event.GuildID,
				ChannelID: event.ChannelID,
				UserID:    event.UserID,
				Type:      events.MessageReactionAddType,
			},
			map[string]interface{}{
				"messageID": messages[0].ID,
				"category":  string(categoryData),
			},
		)
		if err != nil {
			event.Except(err)
			return
		}

		err = discord.React(
			event.Redis(), event.Discord(), messages[0].ChannelID, messages[0].ID, false, "✅",
		)
		if err != nil {
			return
		}
		discord.React(
			event.Redis(), event.Discord(), messages[0].ChannelID, messages[0].ID, false, "❌",
		)

		return
	}

	err = p.db.Delete(category.Roles).Delete(category).Error
	if err != nil {
		event.Except(err)
		return
	}

	event.Respond("roles.category.deleted",
		"category", category,
	)
}

func (p *Plugin) toggleCategory(event *events.Event) {
	if len(event.Fields()) < 4 {
		event.Respond("common.invalid-params")
		return
	}

	toggle := true
	switch event.Fields()[1] {
	case "enable":
		break
	case "disable":
		toggle = false
		break
	default:
		event.Respond("common.invalid-params")
		return
	}

	category, err := p.getCategoryByName(event.Fields()[3], event.GuildID)
	if err != nil {
		event.Except(err)
		return
	}

	if category.Name == "" {
		event.Respond("roles.category.does-not-exist")
		return
	}

	category.Enabled = toggle

	err = p.db.Save(category).Error
	if err != nil {
		event.Except(err)
		return
	}

	event.Respond("roles.category.toggle",
		"category", category,
		"toggle", toggle,
	)
}

func (p *Plugin) toggleCategoryVisibility(event *events.Event) {
	if len(event.Fields()) < 4 {
		event.Respond("common.invalid-params")
		return
	}

	toggle := false
	switch event.Fields()[1] {
	case "show":
		break
	case "hide":
		toggle = true
		break
	default:
		event.Respond("common.invalid-params")
		return
	}

	category, err := p.getCategoryByName(event.Fields()[3], event.GuildID)
	if err != nil {
		event.Except(err)
		return
	}

	if category.Name == "" {
		event.Respond("roles.category.does-not-exist")
		return
	}

	category.Hidden = toggle

	err = p.db.Save(category).Error
	if err != nil {
		event.Except(err)
		return
	}

	event.Respond("roles.category.toggle-visibility",
		"category", category,
		"toggle", toggle,
	)
}
