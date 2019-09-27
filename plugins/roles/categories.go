package roles

import (
	"strconv"

	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) createCategory(event *events.Event) {
	if len(event.Fields()) < 5 {
		event.Respond("common.invalid-params")
		return
	}

	name := event.Fields()[3]
	inputChannel := event.Fields()[4]

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

	channel, err := event.State().ChannelFromMention(event.GuildID, inputChannel)
	if err != nil {
		event.Except(err)
		return
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
		ChannelID: channel.ID,
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
	if len(event.Fields()) < 6 {
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
	inputChannel := event.Fields()[5]

	if name == "" {
		event.Respond("roles.category.no-name")
		return
	}

	channel, err := event.State().ChannelFromMention(event.GuildID, inputChannel)
	if err != nil {
		event.Except(err)
		return
	}

	limit := 0
	if len(event.Fields()) >= 7 {
		limit, err = strconv.Atoi(event.Fields()[6])
		if err != nil {
			event.Respond("roles.category.limit-not-number")
			return
		}
	}

	existingCategory.ChannelID = channel.ID
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

	// TODO: check if this category has roles, if it does. confirm before fully deleting

	err = p.db.Delete(category).Error
	if err != nil {
		event.Except(err)
		return
	}

	event.Respond("roles.category.deleted",
		"category", category,
	)
}
