package roles

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"gitlab.com/Cacophony/go-kit/config"
	"gitlab.com/Cacophony/go-kit/events"
)

const (
	PLUS  = "+"
	MINUS = "-"
)

func (p *Plugin) handleUserRoleRequest(event *events.Event) bool {

	// TODO: for performance reasons the channel lookup needs to use redis

	plusMinus := event.MessageCreate.Content[0:1]
	roleInput := event.MessageCreate.Content[1:]

	if plusMinus != PLUS && plusMinus != MINUS {
		return false
	}

	defaultChannelID, err := config.GuildGetString(event.DB(), event.GuildID, guildRoleChannelKey)
	if err != nil {
		event.Except(err)
		return false
	}

	// check default server role channel first
	if event.ChannelID == defaultChannelID {

		// check if user is adding uncategorized role
		uncategorizedRoles, err := p.getUncategorizedRoles(event.GuildID)
		if err != nil {
			event.Except(err)
			return false
		}
		for _, role := range uncategorizedRoles {
			if role.Match(event.State(), roleInput) {

				if plusMinus == PLUS {
					err = p.assignRole(event, role.ServerRoleID)
				} else {
					err = p.removeRole(event, role.ServerRoleID)
				}
				if err != nil {
					event.Except(err)
					return true
				}
				return true
			}
		}
	}

	return false
}

func (p *Plugin) assignRole(event *events.Event, serverRoleID string) error {
	fmt.Println("assign")

	// check if user already has role
	member, err := event.State().Member(event.GuildID, event.UserID)
	if err != nil {
		return err
	}
	for _, userRole := range member.Roles {
		if userRole == serverRoleID {
			return events.NewUserError(
				event.Translate("roles.role.already-assigned", "userMention", member.Mention()),
			)
		}
	}

	// Assign role
	err = event.Discord().Client.GuildMemberRoleAdd(event.GuildID, event.UserID, serverRoleID)
	if err != nil {
		return err
	}

	role, err := event.State().Role(event.GuildID, serverRoleID)
	if err != nil {
		return err
	}

	_, err = event.Respond("roles.role.assigned", "userMention", member.Mention(), "serverRoleName", role.Name)
	if err != nil {
		return err
	}

	return nil
}

func (p *Plugin) removeRole(event *events.Event, serverRoleID string) error {

	fmt.Println("remove")

	// confirm the user has the role
	member, err := event.State().Member(event.GuildID, event.UserID)
	if err != nil {
		return err
	}
	fmt.Println(serverRoleID)
	spew.Dump(member.Roles)
	var hasRole bool
	for _, userRole := range member.Roles {
		if userRole == serverRoleID {
			hasRole = true
			break
		}
	}

	if !hasRole {
		return events.NewUserError(
			event.Translate("roles.role.not-assigned", "userMention", member.Mention()),
		)
	}

	// Remove role
	err = event.Discord().Client.GuildMemberRoleRemove(event.GuildID, event.UserID, serverRoleID)
	if err != nil {
		return err
	}

	role, err := event.State().Role(event.GuildID, serverRoleID)
	if err != nil {
		return err
	}

	_, err = event.Respond("roles.role.removed-role", "userMention", member.Mention(), "serverRoleName", role.Name)
	if err != nil {
		return err
	}

	return nil
}
