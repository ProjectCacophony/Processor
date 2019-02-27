package whitelist

import (
	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/permissions"
)

// serversPerUserLimit returns the server limit for the current user
// returns < 0 if the limit is unlimited
func serversPerUserLimit(event *events.Event) int {
	if event.Has(permissions.BotOwner) {
		return -1
	}

	return 4
}

func (p *Plugin) extractGuild(discord *discordgo.Session, text string) (*discordgo.Guild, error) {
	if p.snowflakeRegex.MatchString(text) {
		return &discordgo.Guild{
			ID: text,
		}, nil
	}

	var inviteCode string
	parts := p.discordInviteRegex.FindStringSubmatch(text)
	if len(parts) >= 6 {
		inviteCode = parts[5]
	}

	invite, err := discord.Invite(inviteCode)
	if err != nil {
		return nil, err
	}

	return invite.Guild, nil
}
