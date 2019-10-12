package roles

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/config"
	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) displayRoleMessage(event *events.Event) {

	var targetChannel *discordgo.Channel
	var err error

	// find target channel
	if len(event.Fields()) >= 3 {
		targetChannel, err = event.State().ChannelFromMention(event.GuildID, event.Fields()[2])
		if err != nil {
			event.Respond("common.channel-not-found")
			return
		}
	} else {
		targetChannel, err = event.State().Channel(event.ChannelID)
		if err != nil {
			event.Except(err)
			return
		}
	}

	// get categories for channel
	categories, err := p.getCategoryByChannel(targetChannel.ID)
	if err != nil {
		event.Except(err)
		return
	}

	// if channel is also the default role channel, also get uncategorized roles
	channelID, err := config.GuildGetString(event.DB(), event.GuildID, guildRoleChannelKey)
	if err != nil {
		event.Except(err)
		return
	}

	var uncategorizedRoles []*Role
	if targetChannel.ID == channelID {
		uncategorizedRoles, err = p.getUncategorizedRoles(event.GuildID)
		if err != nil {
			event.Except(err)
			return
		}
	}

	var roleNames []string
	for _, role := range uncategorizedRoles {
		roleNames = append(roleNames, fmt.Sprintf("`%s`", role.Name(event.State())))
	}

	var categoryList string
	for _, category := range categories {
		if !category.Enabled || category.Hidden || len(category.Roles) == 0 {
			continue
		}

		var roleNames []string
		for _, role := range category.Roles {
			roleNames = append(roleNames, fmt.Sprintf("`%s`", role.Name(event.State())))
		}

		var limitText string
		if category.Limit > 0 {
			limitText = fmt.Sprintf(" (Limit: %d)", category.Limit)
		}

		categoryList += fmt.Sprintf("\n\n**%s**%s\nRoles: %s", category.Name, limitText, strings.Join(roleNames, ", "))
	}

	_, err = event.Respond("roles.message",
		"uncategorizedRolesList", strings.Join(roleNames, ", "),
		"categoryList", categoryList,
	)
	if err != nil {
		event.Except(err)
		return
	}
}
