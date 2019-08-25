package customcommands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/config"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/permissions"
)

func (p *Plugin) toggleCreatePermission(event *events.Event) {
	if len(event.Fields()) > 3 {
		event.Respond("common.invalid-params")
		return
	}

	// check if changing permissions for a role
	var inputRole *discordgo.Role
	if len(event.Fields()) == 3 {
		var err error
		inputRole, _ = event.State().RoleFromMention(event.GuildID, event.Fields()[2])
		if inputRole == nil {
			inputRole, err = event.State().Role(event.GuildID, event.Fields()[2])
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

		// check if adding or removing
		var isRemoving bool
		for i, roleID := range roleIDs {
			if roleID == inputRole.ID {
				roleIDs = append(roleIDs[:i], roleIDs[i+1:]...)
				isRemoving = true
			}
		}

		if !isRemoving {
			roleIDs = append(roleIDs, inputRole.ID)
		}

		err = config.GuildSetString(p.db, event.GuildID, rolesCreatePermissionToggleKey, strings.Join(roleIDs, ","))
		if err != nil {
			event.Except(err)
			return
		}

		event.Respond("customcommands.permission-toggle", "who", inputRole.Name, "canAdd", !isRemoving)
		return
	}

	current, err := config.GuildGetBool(p.db, event.GuildID, everyoneCreatePermissionKey)
	if err != nil && err.Error() == "invalid Guild ID" {
		event.Except(err)
		return
	}

	err = config.GuildSetBool(p.db, event.GuildID, everyoneCreatePermissionKey, !current)
	if err != nil {
		event.Except(err)
		return
	}

	event.Respond("customcommands.permission-toggle", "who", "everyone", "canAdd", !current)
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
