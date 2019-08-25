package customcommands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/config"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/permissions"
)

func (p *Plugin) viewAllPermissions(event *events.Event) {
	if len(event.Fields()) > 2 {
		event.Respond("common.invalid-params")
		return
	}

	if !event.Has(permissions.DiscordAdministrator) {
		return
	}

	// get use permissions
	cantUseServer, err := config.GuildGetBool(p.db, event.GuildID, denyServerCommandsUsePermissionKey)
	if err != nil && err.Error() == "invalid Guild ID" {
		event.Except(err)
		return
	}
	cantUseUser, err := config.GuildGetBool(p.db, event.GuildID, denyUserCommandsUsePermissionKey)
	if err != nil && err.Error() == "invalid Guild ID" {
		event.Except(err)
		return
	}

	// get create permissions
	canEveryoneCreate, err := config.GuildGetBool(p.db, event.GuildID, everyoneCreatePermissionKey)
	if err != nil && err.Error() == "invalid Guild ID" {
		event.Except(err)
		return
	}
	curRoles, err := config.GuildGetString(p.db, event.GuildID, rolesCreatePermissionToggleKey)
	if err != nil && err.Error() == "invalid Guild ID" {
		event.Except(err)
		return
	}

	var roleIDs []string
	if curRoles != "" {
		roleIDs = strings.Split(curRoles, ",")
	}

	var roleNames []string
	for _, roleID := range roleIDs {
		role, err := event.State().Role(event.GuildID, roleID)
		if err != nil || role == nil {
			continue
		}
		roleNames = append(roleNames, role.Name)
	}

	event.Respond("customcommands.permissions",
		"cantUseServer", cantUseServer,
		"cantUseUser", cantUseUser,
		"canEveryoneCreate", canEveryoneCreate,
		"enabledRoles", strings.Join(roleNames, ", "),
	)

}

func (p *Plugin) toggleCreatePermission(event *events.Event, enable bool) {
	if len(event.Fields()) > 4 {
		event.Respond("common.invalid-params")
		return
	}

	// check if changing permissions for a role
	var inputRole *discordgo.Role
	if len(event.Fields()) == 4 {
		var err error
		inputRole, _ = event.State().RoleFromMention(event.GuildID, event.Fields()[3])
		if inputRole == nil {
			inputRole, err = event.State().Role(event.GuildID, event.Fields()[3])
			if err != nil {
				event.Respond("common.no-role")
				return
			}
		}
	}

	if inputRole != nil {
		curRoles, err := config.GuildGetString(p.db, event.GuildID, rolesCreatePermissionToggleKey)
		if err != nil && err.Error() == "invalid Guild ID" {
			event.Except(err)
			return
		}

		var roleIDs []string
		if curRoles != "" {
			roleIDs = strings.Split(curRoles, ",")
		}

		// if we're removing then find and remove it,
		//  otherwise if adding a role check if it already exists
		var roleExists bool
		for i, roleID := range roleIDs {
			if roleID == inputRole.ID {
				roleExists = true

				if !enable {
					roleIDs = append(roleIDs[:i], roleIDs[i+1:]...)
				}
			}
		}

		if !roleExists && enable {
			roleIDs = append(roleIDs, inputRole.ID)
		}

		err = config.GuildSetString(p.db, event.GuildID, rolesCreatePermissionToggleKey, strings.Join(roleIDs, ","))
		if err != nil {
			event.Except(err)
			return
		}

		event.Respond("customcommands.permission-toggle", "who", inputRole.Name, "canAdd", enable)
		return
	}

	err := config.GuildSetBool(p.db, event.GuildID, everyoneCreatePermissionKey, enable)
	if err != nil {
		event.Except(err)
		return
	}

	event.Respond("customcommands.permission-toggle", "who", "everyone", "canAdd", enable)
}

func (p *Plugin) toggleUsePermission(event *events.Event, disable bool) {
	if len(event.Fields()) > 3 {
		event.Respond("common.invalid-params")
		return
	}
	if len(event.Fields()) == 3 && event.Fields()[2] != "user" {
		event.Respond("common.invalid-params")
		return
	}

	key := denyServerCommandsUsePermissionKey
	var isUserChange bool
	if len(event.Fields()) == 3 && event.Fields()[2] == "user" {
		key = denyUserCommandsUsePermissionKey
		isUserChange = true
	}

	p.logger.Info("here")
	err := config.GuildSetBool(p.db, event.GuildID, key, disable)
	if err != nil {
		event.Except(err)
		return
	}

	event.Respond("customcommands.permission-use-toggle", "level", isUserChange, "cantUse", disable)
}

func (p *Plugin) canEditCommand(event *events.Event) bool {
	if isUserOperation(event) {
		return true
	}

	if event.Has(permissions.DiscordAdministrator) {
		return true
	}

	canEveryone, _ := config.GuildGetBool(p.db, event.GuildID, everyoneCreatePermissionKey)
	if canEveryone {
		return true
	}

	curRoles, _ := config.GuildGetString(p.db, event.GuildID, rolesCreatePermissionToggleKey)
	if curRoles == "" {
		return false
	}

	roleIDs := strings.Split(curRoles, ",")

	member, err := event.State().Member(event.GuildID, event.UserID)
	if err != nil {
		return false
	}
	for _, userRole := range member.Roles {
		for _, allowedrole := range roleIDs {
			if userRole == allowedrole {
				return true
			}
		}
	}

	return false
}

func (p *Plugin) canUseServerCommand(event *events.Event) bool {
	if event.Has(permissions.DiscordAdministrator) {
		return true
	}

	cantUse, _ := config.GuildGetBool(p.db, event.GuildID, denyServerCommandsUsePermissionKey)
	return !cantUse
}

func (p *Plugin) canUseUserCommand(event *events.Event) bool {
	if event.Has(permissions.DiscordAdministrator) {
		return true
	}

	cantUse, _ := config.GuildGetBool(p.db, event.GuildID, denyUserCommandsUsePermissionKey)
	return !cantUse
}
