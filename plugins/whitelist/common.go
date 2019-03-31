package whitelist

import (
	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/discord"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/permissions"
)

// serversPerUserLimit returns the server limit for the current user
// returns < 0 if the limit is unlimited
func serversPerUserLimit(event *events.Event) int {
	if event.Has(permissions.BotOwner) {
		return -1
	}

	return 3
}

func (p *Plugin) extractGuild(session *discord.Session, text string) (*discordgo.Guild, error) {
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

	invite, err := discord.Invite(p.redis, session, inviteCode)
	if err != nil {
		return nil, err
	}

	return invite.Guild, nil
}

func inviteURL(botID string) string {
	return "https://discordapp.com/oauth2/authorize?client_id=" + botID + "&scope=bot&permissions=0x00000000"
}
