package stats

import (
	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/permissions"
)

type findMatches struct {
	Guild   *discordgo.Guild
	User    *discordgo.User
	Channel *discordgo.Channel
	Role    *discordgo.Role
	Emoji   *discordgo.Emoji
}

func (p *Plugin) handleFind(event *events.Event, ID string) {
	var err error
	var targetGuild *discordgo.Guild

	if event.Has(permissions.BotAdmin) && len(event.Fields()) >= 2 {
		targetGuild, err = event.State().Guild(event.Fields()[1])
		if err != nil {
			event.Except(err)
			return
		}
	}

	if targetGuild == nil {
		targetGuild, err = event.State().Guild(event.GuildID)
		if err != nil {
			event.Except(err)
			return
		}
	}

	var matches findMatches

	matches.Guild, _ = event.State().Guild(ID)

	matches.User, _ = event.State().User(ID)

	matches.Channel, _ = event.State().Channel(ID)

	matches.Role, _ = event.State().Role(targetGuild.ID, ID)

	matches.Emoji, _ = event.State().Emoji(targetGuild.ID, ID)

	_, err = event.Respond(
		"stats.find.response",
		"matches", matches,
		"targetGuild", targetGuild,
	)
	event.Except(err)
}
