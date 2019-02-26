package whitelist

import (
	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) whitelistStatus(event *events.Event) {
	servers, err := whitelistFindMany(p.db,
		"whitelisted_by_user_id = ?", event.UserID)
	if err != nil {
		event.Except(err)
		return
	}

	serversEnhanced := make([]enhancedEntry, len(servers))
	for i, server := range servers {

		serversEnhanced[i].Guild, err = p.state.Guild(server.GuildID)
		if err != nil {
			serversEnhanced[i].Guild = &discordgo.Guild{
				ID:   server.GuildID,
				Name: "N/A",
			}
		}

		serversEnhanced[i].At = server.UpdatedAt
	}

	_, err = event.Respond("whitelist.status.message",
		"servers", serversEnhanced,
	)
	event.Except(err)
}
