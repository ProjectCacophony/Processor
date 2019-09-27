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
		event.Respond("roles.role.role-not-found-on-server")
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
		if existingCategory.Name != "" {
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

	err = p.db.Save(role).Error
	if err != nil {
		event.Except(err)
		return
	}

	event.Respond("roles.role.creation",
		"roleName", roleName,
	)
}

// func (p *Plugin) updateCategory(event *events.Event) {
// 	if len(event.Fields()) < 6 {
// 		event.Respond("common.invalid-params")
// 		return
// 	}

// 	currentName := event.Fields()[3]
// 	existingCategory, err := p.getCategoryByName(currentName, event.GuildID)
// 	if err != nil {
// 		event.Except(err)
// 		return
// 	}
// 	if existingCategory.Name == "" {
// 		event.Respond("roles.category.does-not-exist")
// 		return
// 	}

// 	name := event.Fields()[4]
// 	inputChannel := event.Fields()[5]

// 	if name == "" {
// 		event.Respond("roles.category.no-name")
// 		return
// 	}

// 	channel, err := event.State().ChannelFromMention(event.GuildID, inputChannel)
// 	if err != nil {
// 		event.Except(err)
// 		return
// 	}

// 	limit := 0
// 	if len(event.Fields()) >= 7 {
// 		limit, err = strconv.Atoi(event.Fields()[6])
// 		if err != nil {
// 			event.Respond("roles.category.limit-not-number")
// 			return
// 		}
// 	}

// 	existingCategory.ChannelID = channel.ID
// 	existingCategory.Name = name
// 	existingCategory.Limit = limit

// 	err = p.db.Save(existingCategory).Error
// 	if err != nil {
// 		event.Except(err)
// 		return
// 	}

// 	event.Respond("roles.category.update",
// 		"category", existingCategory,
// 	)
// }

// func (p *Plugin) deleteCategory(event *events.Event) {
// 	if len(event.Fields()) < 4 {
// 		event.Respond("common.invalid-params")
// 		return
// 	}

// 	category, err := p.getCategoryByName(event.Fields()[3], event.GuildID)
// 	if err != nil {
// 		event.Except(err)
// 		return
// 	}

// 	if category.Name == "" {
// 		event.Respond("roles.category.does-not-exist")
// 		return
// 	}

// 	// TODO: check if this category has roles, if it does. confirm before fully deleting

// 	err = p.db.Delete(category).Error
// 	if err != nil {
// 		event.Except(err)
// 		return
// 	}

// 	event.Respond("roles.category.deleted",
// 		"category", category,
// 	)
// }
