package roles

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/config"
	"gitlab.com/Cacophony/go-kit/events"
)

const (
	PLUS  = "+"
	MINUS = "-"
)

func (p *Plugin) handleUserRoleRequest(event *events.Event) bool {

	if event.MessageCreate == nil {
		return false
	}

	// explicitly not checking error here
	defaultChannelID, _ := config.GuildGetString(event.DB(), event.GuildID, guildRoleChannelKey)

	// check if default server role channel first
	if defaultChannelID != "" && event.ChannelID == defaultChannelID {
		go p.deleteWithDelay(event, event.MessageID)

		if len(event.MessageCreate.Content) < 2 {
			return false
		}
		plusMinus := event.MessageCreate.Content[0:1]
		roleInput := strings.TrimSpace(event.MessageCreate.Content[1:])

		if plusMinus != PLUS && plusMinus != MINUS {
			return false
		}

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

	// TODO: for performance reasons the channel lookup needs to use redis

	// get categories setup for the given channel
	categories, err := p.getCategoryByChannel(event.ChannelID)
	if err != nil {
		event.Except(err)
		return false
	}
	if len(categories) == 0 {
		return false
	}
	go p.deleteWithDelay(event, event.MessageID)

	if len(event.MessageCreate.Content) < 2 {
		return false
	}
	plusMinus := event.MessageCreate.Content[0:1]
	roleInput := strings.TrimSpace(event.MessageCreate.Content[1:])

	if plusMinus != PLUS && plusMinus != MINUS {
		return false
	}

	var member *discordgo.Member
	for _, category := range categories {

		if !category.Enabled {
			continue
		}

		for _, role := range category.Roles {
			if role.Match(event.State(), roleInput) {

				if member == nil {
					member, err = event.State().Member(event.GuildID, event.UserID)
					if err != nil {
						event.Except(err)
						return false
					}
				}

				if plusMinus == PLUS {

					if p.isOverRoleLimit(member, category) {
						msgs, err := event.Respond("roles.role.at-category-limit", "userMention", member.Mention())
						if err != nil {
							return false
						}

						go p.deleteWithDelay(event, msgs[0].ID)
						return true
					}

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

func (p *Plugin) isOverRoleLimit(member *discordgo.Member, category *Category) bool {
	if category.Limit == 0 {
		return false
	}

	if len(member.Roles) < category.Limit {
		return false
	}

	var hasRoleCount int
	for _, role := range category.Roles {
		for _, userRoleID := range member.Roles {
			if role.ServerRoleID == userRoleID {
				hasRoleCount++
			}
		}
	}

	if hasRoleCount >= category.Limit {
		return true
	}

	return false
}

func (p *Plugin) assignRole(event *events.Event, serverRoleID string) error {

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

	msgs, err := event.Respond("roles.role.assigned", "userMention", member.Mention(), "serverRoleName", role.Name)
	if err != nil {
		return err
	}
	go p.deleteWithDelay(event, msgs[0].ID)

	return nil
}

func (p *Plugin) removeRole(event *events.Event, serverRoleID string) error {

	// confirm the user has the role
	member, err := event.State().Member(event.GuildID, event.UserID)
	if err != nil {
		return err
	}
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

	msgs, err := event.Respond("roles.role.removed-role", "userMention", member.Mention(), "serverRoleName", role.Name)
	if err != nil {
		return err
	}

	go p.deleteWithDelay(event, msgs[0].ID)

	return nil
}
