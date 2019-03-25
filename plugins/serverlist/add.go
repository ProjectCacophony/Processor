package serverlist

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/regexp"
)

// TODO: match category by direct linking the channel,
//       also allow sub linking the channel if category

// serverlist add discord.gg/? "Join my awesome server!" "girl group; boy group"
func (p *Plugin) handleAdd(event *events.Event) {
	if len(event.Fields()) < 5 {
		event.Respond("serverlist.add.too-few") // nolint: errcheck
		return
	}

	fields := event.Fields()[2:]

	var err error
	var invite *discordgo.Invite
	parts := regexp.DiscordInviteRegexp.FindStringSubmatch(fields[0])
	if len(parts) >= 6 {
		invite, err = event.Discord().Client.InviteWithCounts(parts[5])
		if err != nil {
			event.Except(err)
			return
		}
	}

	if invite == nil || invite.Code == "" || invite.Guild.ID == "" {
		event.Respond("serverlist.add.invalid-invite") // nolint: errcheck
		return
	}

	description := fields[1]
	if len(description) > descriptionCharacterLimit {
		event.Respond("serverlist.add.description-too-long", // nolint: errcheck
			"limit", descriptionCharacterLimit,
		)
		return
	}

	allCategories, err := categoriesFindMany(p.db, "bot_id = ?", event.BotUserID)
	if err != nil {
		event.Except(err)
		return
	}

	var categoryIDs []uint
	for _, categoryName := range strings.Split(fields[2], ";") {
		categoryName = strings.ToLower(strings.TrimSpace(categoryName))

		for _, allCategory := range allCategories {
			for _, keyword := range allCategory.Keywords {
				if keyword != categoryName {
					continue
				}

				if uintSliceContains(allCategory.ID, categoryIDs) {
					continue
				}

				categoryIDs = append(categoryIDs, allCategory.ID)
			}
		}
	}

	err = serverAdd(
		p.db,
		[]string{invite.Guild.Name},
		description,
		invite.Code,
		invite.Guild.ID,
		[]string{event.UserID},
		categoryIDs,
		invite.ApproximateMemberCount,
		event.BotUserID,
	)
	if err != nil {
		event.Except(err)
		return
	}

	_, err = event.Respond("serverlist.add.content",
		"name", invite.Guild.Name,
	)
	event.Except(err)
}

func uintSliceContains(n uint, list []uint) bool {
	for _, i := range list {
		if n != i {
			continue
		}

		return true
	}

	return false
}
