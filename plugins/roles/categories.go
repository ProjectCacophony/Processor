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

	existingCategory, err := p.getCategoryByName(name)
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

func (p *Plugin) getCategoryByName(name string) (*Category, error) {
	var category Category
	err := p.db.First(&category, "name = ?", name).Error
	if err != nil && !strings.Contains(err.Error(), "record not found") {
		return nil, err
	}
	return &category, nil
}
