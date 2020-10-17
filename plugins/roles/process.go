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

	if event.MessageCreate == nil || event.MessageCreate.Author == nil || event.MessageCreate.Author.Bot {
		return false
	}

	// check if the message was sent in a role channel
	inRoleChannel := false
	channels := p.getCachedRoleChannels(event.GuildID)
	if len(channels) > 0 {
		for _, channel := range channels {
			if channel == event.ChannelID {
				inRoleChannel = true
				break
			}
		}
	}
	if !inRoleChannel {
		return false
	}

	// remove users message
	go p.deleteWithDelay(event, event.MessageID)

	// check plus or minus
	if len(event.MessageCreate.Content) < 2 {
		return false
	}
	plusMinus := event.MessageCreate.Content[0:1]
	if plusMinus != PLUS && plusMinus != MINUS {
		return false
	}

	// get all the roles setup for the channel or categories using the channel
	uncategorizedRoles, err := p.getUncategorizedRoles(event.GuildID)
	if err != nil {
		event.Except(err)
		return false
	}
	categories, err := p.getCategoryByChannel(event.ChannelID)
	if err != nil {
		event.Except(err)
		return false
	}

	allRoles := append([]*Role{}, uncategorizedRoles...)

	for _, category := range categories {
		for _, role := range category.Roles {
			role := role
			allRoles = append(allRoles, &role)
		}
	}

	requests, err := p.parseRoleRequestMessage(event, strings.TrimSpace(event.MessageCreate.Content), allRoles)
	if err != nil {
		event.Except(err)
		return true
	}

	// explicitly not checking error here
	defaultChannelID, _ := config.GuildGetString(event.DB(), event.GuildID, guildRoleChannelKey)

	for role, plusMinus := range requests {

		// check if default server role channel first
		if defaultChannelID != "" && event.ChannelID == defaultChannelID {

			for _, urole := range uncategorizedRoles {
				if urole.ServerRoleID == role.ServerRoleID {

					if string(plusMinus) == PLUS {
						err = p.assignRole(event, role.ServerRoleID)
					} else {
						err = p.removeRole(event, role.ServerRoleID)
					}
					if err != nil {
						event.Except(err)
						break
					}
					break
				}
			}
		}

		if len(categories) == 0 {
			continue
		}

		var member *discordgo.Member
		for _, category := range categories {

			if !category.Enabled {
				continue
			}

			for _, crole := range category.Roles {
				if crole.ServerRoleID == role.ServerRoleID {

					if member == nil {
						member, err = event.State().Member(event.GuildID, event.UserID)
						if err != nil {
							event.Except(err)
							return false
						}
					}

					if string(plusMinus) == PLUS {

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
						break
					}
					break
				}
			}
		}
	}

	return false
}

func (p *Plugin) parseRoleRequestMessage(event *events.Event, userMsg string, roles []*Role) (map[*Role]rune, error) {
	requestsMap := make(map[*Role]rune)
	ignoreChars := 0
	skipSpace := false

	for i, str := range userMsg {
		// might be spaces between the +/- and the previous assigned role name, need
		// to skip spaces before we start counting characters
		if skipSpace {
			if str != ' ' {
				skipSpace = false
			} else {
				continue
			}
		}

		// if a role was found, skip the amount of characters that match the role name length
		if ignoreChars != 0 {
			ignoreChars--
			continue
		}

		if str == '+' || str == '-' {
			remainingMsg := strings.TrimSpace(userMsg[i+1:])
			foundRole := false
		RoleLoop:
			for _, role := range roles {
				roleName := role.Name(p.state)
				if len(roleName) != 0 && strings.HasPrefix(remainingMsg, roleName) {
					ignoreChars = len(roleName)
					requestsMap[role] = str
					foundRole = true
					skipSpace = true
					break RoleLoop
				}

				for _, alias := range role.Aliases {
					if len(alias) != 0 && strings.HasPrefix(remainingMsg, alias) {
						ignoreChars = len(alias)
						requestsMap[role] = str
						foundRole = true
						skipSpace = true
						break RoleLoop
					}
				}
			}

			if !foundRole {
				member, err := p.state.Member(event.GuildID, event.UserID)
				if err != nil {
					return nil, err
				}
				unfoundRole := strings.Split(remainingMsg, " ")[0]
				return nil, events.NewUserError(
					event.Translate("roles.role.rolename-not-found", "userMention", member.Mention(), "roleName", unfoundRole),
				)
			}
		}
	}

	return requestsMap, nil
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

	return hasRoleCount >= category.Limit
}

func (p *Plugin) assignRole(event *events.Event, serverRoleID string) error {

	// check if user already has role
	member, err := event.State().Member(event.GuildID, event.UserID)
	if err != nil {
		return err
	}

	role, err := event.State().Role(event.GuildID, serverRoleID)
	if err != nil {
		return err
	}

	for _, userRole := range member.Roles {
		if userRole == serverRoleID {
			msgs, err := event.Respond("roles.role.already-assigned", "userMention", member.Mention(), "serverRoleName", role.Name)
			if len(msgs) > 0 {
				go p.deleteWithDelay(event, msgs[0].ID)
			}
			return err
		}
	}

	// Assign role
	err = event.Discord().Client.GuildMemberRoleAdd(event.GuildID, event.UserID, serverRoleID)
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

	role, err := event.State().Role(event.GuildID, serverRoleID)
	if err != nil {
		return err
	}

	if !hasRole {
		msgs, err := event.Respond("roles.role.not-assigned", "userMention", member.Mention(), "serverRoleName", role.Name)
		if len(msgs) > 0 {
			go p.deleteWithDelay(event, msgs[0].ID)
		}
		return err
	}

	// Remove role
	err = event.Discord().Client.GuildMemberRoleRemove(event.GuildID, event.UserID, serverRoleID)
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
