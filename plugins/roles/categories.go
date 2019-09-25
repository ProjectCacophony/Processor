package roles

import (
	"strconv"
	"strings"

	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) createCategory(event *events.Event) {
	if len(event.Fields()) < 6 {
		event.Respond("common.invalid-params")
		return
	}

	name := event.Fields()[3]
	message := event.Fields()[4]
	inputChannel := event.Fields()[5]

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
	if len(event.Fields()) >= 7 {
		limit, err = strconv.Atoi(event.Fields()[6])
		if err != nil {
			event.Respond("roles.category.limit-not-number")
			return
		}
	}

	pool := ""
	if len(event.Fields()) >= 8 {
		pool = event.Fields()[7]
	}

	category := &Category{
		GuildID:   event.GuildID,
		ChannelID: channel.ID,
		Name:      name,
		Message:   message,
		Limit:     limit,
		Enabled:   true,
		Pool:      pool,
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
	if len(event.Fields()) < 7 {
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
	message := event.Fields()[5]
	inputChannel := event.Fields()[6]

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
	if len(event.Fields()) >= 8 {
		limit, err = strconv.Atoi(event.Fields()[7])
		if err != nil {
			event.Respond("roles.category.limit-not-number")
			return
		}
	}

	pool := ""
	if len(event.Fields()) >= 9 {
		pool = event.Fields()[8]
	}

	existingCategory.ChannelID = channel.ID
	existingCategory.Name = name
	existingCategory.Message = message
	existingCategory.Pool = pool
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

func (p *Plugin) getCategoryByName(name string, guildId string) (*Category, error) {
	var category Category
	err := p.db.First(&category, "name = ? and guild_id = ?", name, guildId).Error
	if err != nil && !strings.Contains(err.Error(), "record not found") {
		return nil, err
	}
	return &category, nil
}
