package roles

import (
	"strings"

	"gitlab.com/Cacophony/go-kit/events"
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

	serverRole, err := p.getServerRoleByNameOrID(serverRoleID, event.GuildID)
	if err != nil {
		if err.Error() == ServerRoleNotFound || err.Error() == MultipleServerRolesWithName {
			event.Respond(err.Error())
		} else {
			event.Except(err)
		}
		return
	}
	serverRoleID = serverRole.ID
	roleName := serverRole.Name

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
		ServerRoleID: serverRoleID,
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

	err = p.db.Save(role).Error
	if err != nil {
		event.Except(err)
		return
	}

	event.Respond("roles.role.creation",
		"roleName", roleName,
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
		if existingCategory.Name == "" {
			event.Respond("roles.category.does-not-exist")
			return
		}

		categoryID = existingCategory.ID
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
