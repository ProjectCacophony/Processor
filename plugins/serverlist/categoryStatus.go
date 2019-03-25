package serverlist

import (
	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/events"
)

type enhancedCategory struct {
	Category
	By *discordgo.User
}

func (p *Plugin) handleCategoryStatus(event *events.Event) {
	categories, err := categoriesFindMany(p.db, "guild_id = ?", event.GuildID)
	if err != nil {
		event.Except(err)
		return
	}

	enhancedCategories := make([]enhancedCategory, len(categories))
	for i, category := range categories {
		enhancedCategories[i].Category = category

		user, err := p.state.User(category.AddedBy)
		if err != nil {
			user = &discordgo.User{
				Username: "N/A",
				ID:       category.AddedBy,
			}
		}
		enhancedCategories[i].By = user
	}

	_, err = event.Respond("serverlist.category-status.content",
		"categories", enhancedCategories,
	)
	event.Except(err)
}
