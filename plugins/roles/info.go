package roles

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) displayRoleInfo(event *events.Event) {

	categories, err := p.getAllCategories(event.GuildID)
	if err != nil {
		event.Except(err)
		return
	}

	roles, err := p.getUncategorizedRoles(event.GuildID)
	if err != nil {
		event.Except(err)
		return
	}

	if len(categories) == 0 && len(roles) == 0 {
		event.Respond("roles.category.no-categories-or-roles")
		return
	}

	outputText := ""

	categoriesText := ""

	if len(roles) > 0 {
		roleText := ""

		for _, role := range roles {
			roleText += p.formatRoleOutput(role, event.GuildID)
		}

		categoryText := fmt.Sprintf("**%s** \n%s\n",
			"Uncategorized Roles",
			roleText,
		)
		categoriesText += categoryText

	}

	for _, cat := range categories {
		spew.Dump(cat)

		status := "Enabled"
		if !cat.Enabled {
			status = "Disabled"
		}

		channelName := "*Unknown*"
		channel, err := event.State().Channel(cat.ChannelID)
		if err == nil {
			channelName = channel.Name
		}

		limitText := "No Limit"
		if cat.Limit > 0 {
			limitText = fmt.Sprintf("Limit: %d", cat.Limit)
		}

		roleText := "No Roles"
		if len(cat.Roles) > 0 {
			roleText = ""

			for _, role := range cat.Roles {
				roleText += p.formatRoleOutput(&role, event.GuildID)
			}
		}

		categoryText := fmt.Sprintf("**%s** (%s, #%s, %s)\n%s\n",
			cat.Name,
			limitText,
			channelName,
			status,
			roleText,
		)
		categoriesText += categoryText
	}

	outputText += "\n\n" + categoriesText

	event.Respond(outputText)
}

func (p *Plugin) formatRoleOutput(role *Role, guildID string) string {
	roleText := ""
	serverRole, err := p.state.Role(guildID, role.ServerRoleID)
	if err != nil {
		return ""
	}

	roleText += fmt.Sprintf("\t**%s** (%s) ", serverRole.Name, role.ServerRoleID)
	if role.PrintName != "" {
		roleText += fmt.Sprintf("__Print__=%s ", role.PrintName)
	}
	if len(role.Aliases) > 0 && role.Aliases[0] != "" {
		roleText += fmt.Sprintf("__Aliases__=%s ", role.Aliases)
	}

	roleText += "\n"

	return roleText
}
