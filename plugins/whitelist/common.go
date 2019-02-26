package whitelist

import "github.com/bwmarrin/discordgo"

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
