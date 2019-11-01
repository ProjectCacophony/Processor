package roles

import (
	"encoding/json"
	"strings"

	"gitlab.com/Cacophony/go-kit/discord"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/permissions"
)

func (p *Plugin) createRole(event *events.Event) {
	if len(event.Fields()) < 4 {
		event.Respond("common.invalid-params")
		return
	}

	serverRoleID := event.Fields()[3]

	if serverRoleID == "" {
		event.Respond("roles.role.no-name")
		return
	}

	existingRole, err := p.getRoleByServerRoleID(serverRoleID, event.GuildID)
	if err != nil {
		event.Except(err)
		return
	}
	if existingRole.ServerRoleID != "" {
		event.Respond("roles.role.role-already-setup")
		return
	}

	var printName string
	if len(event.Fields()) >= 5 {
		printName = event.Fields()[4]
	}

	var aliases []string
	if len(event.Fields()) >= 6 {
		aliases = strings.Split(event.Fields()[5], ",")
		for i, v := range aliases {
			aliases[i] = strings.TrimSpace(v)
		}
	}

	var categoryID uint
	if len(event.Fields()) >= 7 {

		existingCategory, err := p.getCategoryByName(event.Fields()[6], event.GuildID)
		if err != nil {
			event.Except(err)
			return
		}
		if existingCategory.Name == "" {
			event.Respond("roles.category.does-not-exist")
			return
		}

		categoryID = existingCategory.ID
	}

	role := &Role{
		ServerRoleID: "",
		CategoryID:   categoryID,
		PrintName:    printName,
		Aliases:      aliases,
		GuildID:      event.GuildID,
		Enabled:      true,
	}

	if badAlias, ok := p.validateRolePrintAndAliases(role); !ok {
		event.Respond("roles.role.alias-already-exists", "aliasName", badAlias)
		return
	}

	serverRole, err := p.getServerRoleByNameOrID(serverRoleID, event.GuildID)
	if err != nil {

		if err.Error() == ServerRoleNotFound {
			if permissions.DiscordManageRoles.Match(
				event.State(),
				event.DB(),
				event.BotUserID,
				event.ChannelID,
				false,
			) {
				messages, err := event.Respond("roles.roles.ask-role-create")
				if err != nil {
					event.Except(err)
					return
				}

				roleData, err := json.Marshal(&role)
				if err != nil {
					event.Except(err)
					return
				}

				err = event.Questionnaire().Register(
					confirmCreateRoleKey,
					events.QuestionnaireFilter{
						GuildID:   event.GuildID,
						ChannelID: event.ChannelID,
						UserID:    event.UserID,
						Type:      events.MessageReactionAddType,
					},
					map[string]interface{}{
						"messageID": messages[0].ID,
						"role":      string(roleData),
						"roleName":  serverRoleID,
					},
				)
				if err != nil {
					event.Except(err)
					return
				}

				err = discord.React(
					event.Redis(), event.Discord(), messages[0].ChannelID, messages[0].ID, false, "✅",
				)
				if err != nil {
					return
				}
				discord.React(
					event.Redis(), event.Discord(), messages[0].ChannelID, messages[0].ID, false, "❌",
				)

				return
			} else {
				event.Respond(err.Error())
			}
		} else if err.Error() == MultipleServerRolesWithName {
			event.Respond(err.Error())
		} else {
			event.Except(err)
		}
		return
	}
	role.ServerRoleID = serverRole.ID

	err = p.db.Save(role).Error
	if err != nil {
		event.Except(err)
		return
	}

	event.Respond("roles.role.creation",
		"roleName", serverRole.Name,
	)
}

func (p *Plugin) updateRole(event *events.Event) {
	if len(event.Fields()) < 4 {
		event.Respond("common.invalid-params")
		return
	}

	roleID := event.Fields()[3]
	if roleID == "" {
		event.Respond("roles.role.no-name")
		return
	}

	serverRole, err := p.getServerRoleByNameOrID(roleID, event.GuildID)
	if err != nil {
		if err.Error() == ServerRoleNotFound || err.Error() == MultipleServerRolesWithName {
			event.Respond(err.Error())
		} else {
			event.Except(err)
		}
		return
	}
	roleName := serverRole.Name

	existingRole, err := p.getRoleByServerRoleID(serverRole.ID, event.GuildID)
	if err != nil {
		event.Except(err)
		return
	}
	if existingRole.ServerRoleID == "" {
		event.Respond("roles.role.role-not-found")
		return
	}

	var printName string
	if len(event.Fields()) >= 5 {
		printName = event.Fields()[4]
	}

	var aliases []string
	if len(event.Fields()) >= 6 {
		aliases = strings.Split(event.Fields()[5], ",")
		for i, v := range aliases {
			aliases[i] = strings.TrimSpace(v)
		}
	}

	var categoryID uint
	if len(event.Fields()) >= 7 {

		existingCategory, err := p.getCategoryByName(event.Fields()[6], event.GuildID)
		if err != nil {
			event.Except(err)
			return
		}
		if existingCategory.Name != "" {
			categoryID = existingCategory.ID
		}

	}

	existingRole.CategoryID = categoryID
	existingRole.PrintName = printName
	existingRole.Aliases = aliases

	if badAlias, ok := p.validateRolePrintAndAliases(existingRole); !ok {
		event.Respond("roles.role.alias-already-exists", "aliasName", badAlias)
		return
	}

	err = p.db.Save(existingRole).Error
	if err != nil {
		event.Except(err)
		return
	}

	event.Respond("roles.role.update",
		"roleName", roleName,
	)
}

func (p *Plugin) deleteRole(event *events.Event) {
	if len(event.Fields()) < 4 {
		event.Respond("common.invalid-params")
		return
	}

	serverRole, err := p.getServerRoleByNameOrID(event.Fields()[3], event.GuildID)
	if err != nil {
		if err.Error() == ServerRoleNotFound || err.Error() == MultipleServerRolesWithName {
			event.Respond(err.Error())
		} else {
			event.Except(err)
		}
		return
	}

	existingRole, err := p.getRoleByServerRoleID(serverRole.ID, event.GuildID)
	if err != nil {
		event.Except(err)
		return
	}
	if existingRole.ServerRoleID == "" {
		event.Respond("roles.role.role-not-found")
		return
	}

	err = p.db.Delete(existingRole).Error
	if err != nil {
		event.Except(err)
		return
	}

	event.Respond("roles.role.deleted",
		"roleName", serverRole.Name,
	)
}

func (p *Plugin) validateRolePrintAndAliases(role *Role) (string, bool) {
	var values []string

	if role.PrintName != "" {
		values = append(values, role.PrintName)
	}
	for _, alias := range role.Aliases {
		if alias != "" {
			values = append(values, alias)
		}
	}

	if len(values) == 0 {
		return "", true
	}

	for index, value := range values {
		for i, v := range values {
			if index != i && v == value {
				return v, false
			}
		}
	}

	roles, err := p.getAllRoles(role.GuildID)
	if err != nil {
		return "", false
	}

	if len(roles) == 0 {
		return "", true
	}

	for _, r := range roles {
		for _, v := range values {
			if v == r.PrintName {
				return r.PrintName, false
			}

			for _, alias := range r.Aliases {
				if v == alias {
					return alias, false
				}
			}
		}
	}

	return "", true
}
