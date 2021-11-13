package stats

import (
	"regexp"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/permissions"
)

const searchLimit = 10

func (p *Plugin) handleSearch(event *events.Event) {
	if len(event.Fields()) < 3 {
		event.Respond("stats.search.no-search-value")
		return
	}

	searchExp, err := regexp.Compile(event.Fields()[2])
	if err != nil {
		event.Respond("stats.search.search-value-invalid")
		return
	}

	var targetGuild *discordgo.Guild
	if event.Has(permissions.BotAdmin) && len(event.Fields()) >= 4 {
		targetGuild, err = event.State().Guild(event.Fields()[3])
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

	members, err := event.State().GuildMembers(targetGuild.ID)
	if err != nil {
		event.Except(err)
		return
	}

	result := make([]*discordgo.Member, 0, searchLimit)

	for _, memberID := range members {
		member, err := event.State().Member(targetGuild.ID, memberID)
		if err != nil {
			event.Except(err)
			return
		}

		if !(searchExp.MatchString(member.User.Username) ||
			searchExp.MatchString(member.User.String()) ||
			(member.Nick != "" && searchExp.MatchString(member.Nick))) {
			continue
		}

		result = append(result, member)

		if len(result) >= searchLimit {
			break
		}
	}

	event.Respond("stats.search.content", "searchText", searchExp.String(), "members", result)
}
