package roles

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/config"
	"gitlab.com/Cacophony/go-kit/discord"
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

	reactions := make([]string, 0)
	roleNames := make([]string, len(uncategorizedRoles))
	for i, role := range uncategorizedRoles {
		if role.Emoji == "" {
			roleNames[i] = fmt.Sprintf("`%s`", role.Name(event.State()))
		} else {
			reactions = append(reactions, role.Emoji)
			roleNames[i] = fmt.Sprintf("`%s` %s", role.Name(event.State()), role.Emoji)
		}
	}

	var categoryList string
	for _, category := range categories {
		if !category.Enabled || category.Hidden || len(category.Roles) == 0 {
			continue
		}

		var roleNames []string
		for _, role := range category.Roles {
			if role.Emoji == "" {
				roleNames = append(roleNames, fmt.Sprintf("`%s`", role.Name(event.State())))
			} else {
				reactions = append(reactions, role.Emoji)
				roleNames = append(roleNames, fmt.Sprintf("`%s` %s", role.Name(event.State()), role.Emoji))
			}
		}

		var limitText string
		if category.Limit > 0 {
			limitText = fmt.Sprintf(" (Limit: %d)", category.Limit)
		}

		categoryList += fmt.Sprintf("\n\n**%s**%s\nRoles: %s", category.Name, limitText, strings.Join(roleNames, ", "))
	}

	discord.Delete(event.Redis(), event.Discord(), event.ChannelID, event.MessageID, event.DM())
	msg, err := event.Send(targetChannel.ID, "roles.message",
		"uncategorizedRolesList", strings.Join(roleNames, ", "),
		"categoryList", categoryList,
	)
	if err != nil {
		event.Except(err)
		return
	}

	for _, reaction := range reactions {
		discord.React(event.Redis(), event.Discord(), msg[0].ChannelID, msg[0].ID, false, reaction)
	}
}
