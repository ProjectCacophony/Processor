package roles

import (
	"fmt"

	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) displayRoleInfo(event *events.Event) {

	categories, err := p.getAllCategories(event.GuildID)
	if err != nil {
		event.Except(err)
		return
	}

	if len(categories) == 0 {
		event.Respond("roles.category.no-categories")
		return
	}

	// spew.Dump(categories)

	outputText := ""

	categoriesText := "**Category Info:**\n"
	for _, cat := range categories {

		status := "Enabled"
		if !cat.Enabled {
			status = "Disabled"
		}

		channelName := "*Unknown*"
		channel, err := event.State().Channel(cat.ChannelID)
		if err == nil {
			channelName = channel.Name
		}

		limitText := "No Limit"
		if cat.Limit > 0 {
			limitText = fmt.Sprintf("Limit: %d", cat.Limit)
		}

		categoryText := fmt.Sprintf("**%s** (%s, #%s, %s)\n\t%s\n\n",
			cat.Name,
			limitText,
			channelName,
			status,
			"TODO: add category roles here...",
		)
		categoriesText += categoryText
	}

	outputText += "\n\n" + categoriesText

	event.Respond(outputText)
}
