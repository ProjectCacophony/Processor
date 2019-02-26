package whitelist

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/events"
)

type enhancedEntry struct {
	By    *discordgo.User
	Guild *discordgo.Guild
	At    time.Time
}

func (p *Plugin) whitelistList(event *events.Event) {
	whitelistAllServers, err := whitelistAll(p.db)
	if err != nil {
		event.Except(err)
		return
	}

	blacklistAllServers, err := blacklistAll(p.db)
	if err != nil {
		event.Except(err)
		return
	}

	whitelistAllServersEnhanced := make([]enhancedEntry, len(whitelistAllServers))
	for i, server := range whitelistAllServers {

		whitelistAllServersEnhanced[i].Guild, err = p.state.Guild(server.GuildID)
		if err != nil {
			whitelistAllServersEnhanced[i].Guild = &discordgo.Guild{
				ID:   server.GuildID,
				Name: "N/A",
			}
		}

		whitelistAllServersEnhanced[i].By, err = p.state.User(server.WhitelistedByUserID)
		if err != nil {
			whitelistAllServersEnhanced[i].By = &discordgo.User{
				ID:       server.WhitelistedByUserID,
				Username: "N/A",
			}
		}

		whitelistAllServersEnhanced[i].At = server.UpdatedAt
	}

	blacklistAllServersEnhanced := make([]enhancedEntry, len(blacklistAllServers))
	for i, server := range blacklistAllServers {

		blacklistAllServersEnhanced[i].Guild, err = p.state.Guild(server.GuildID)
		if err != nil {
			blacklistAllServersEnhanced[i].Guild = &discordgo.Guild{
				ID:   server.GuildID,
				Name: "N/A",
			}
		}

		blacklistAllServersEnhanced[i].By, err = p.state.User(server.BlacklistedByUserID)
		if err != nil {
			blacklistAllServersEnhanced[i].By = &discordgo.User{
				ID:       server.BlacklistedByUserID,
				Username: "N/A",
			}
		}

		whitelistAllServersEnhanced[i].At = server.UpdatedAt
	}

	_, err = event.Respond("whitelist.list.message",
		"whitelisted", whitelistAllServersEnhanced,
		"blacklisted", blacklistAllServersEnhanced,
	)
	event.Except(err)
}
