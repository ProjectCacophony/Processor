package gall

import (
	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/events"
)

type enhancedEntry struct {
	Entry
	By    *discordgo.User
	Posts int
}

func (p *Plugin) status(event *events.Event) {
	entries, err := entryFindMany(p.db, "guild_id = ?", event.GuildID)
	if err != nil {
		event.Except(err)
		return
	}

	enhancedEntries := make([]enhancedEntry, len(entries))
	for i, entry := range entries {
		enhancedEntries[i].Entry = entry

		user, err := p.state.User(entry.AddedBy)
		if err != nil {
			user = &discordgo.User{
				Username: "N/A",
				ID:       entry.AddedBy,
			}
		}
		enhancedEntries[i].By = user

		enhancedEntries[i].Posts, _ = countPosts(p.db, entry.ID)
	}

	_, err = event.Respond("gall.status.message",
		"entries", enhancedEntries,
		"limit", feedsPerGuildLimit(event),
	)
	event.Except(err)
}
