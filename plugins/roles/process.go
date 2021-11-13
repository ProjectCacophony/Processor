package roles

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/config"
	"gitlab.com/Cacophony/go-kit/discord"
	"gitlab.com/Cacophony/go-kit/discord/emoji"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/permissions"
)

const (
	PLUS  = "+"
	MINUS = "-"
)

func (p *Plugin) handleUserRoleReactionRequest(event *events.Event) bool {
	if event.BotUserID == event.MessageReactionAdd.UserID ||
		event.MessageReactionAdd == nil ||
		event.MessageReactionAdd.Emoji.Name == "" ||
		!p.isInRoleChannel(event) {
		return false
	}

	snowflake := fmt.Sprintf(":%s:%s", event.MessageReactionAdd.Emoji.Name, event.MessageReactionAdd.Emoji.ID)
	if event.MessageReactionAdd.Emoji.Animated {
		snowflake = "a" + snowflake
	}

	go discord.RemoveReact(
		event.Redis(),
		event.Discord(),
		event.MessageReactionAdd.ChannelID,
		event.MessageReactionAdd.MessageID,
		event.UserID,
		false,
		snowflake,
	)

	allRoles, err := p.getAllRoles(event.GuildID)
	if err != nil {
		event.Except(err)
		return true
	}

	var selectedRole *Role
	for _, role := range allRoles {
		if role.Emoji == "" {
			continue
		}

		if snowflake == emoji.GetWithout(role.Emoji) || snowflake == fmt.Sprintf(":%s:", role.Emoji) {
			selectedRole = role
		}
	}

	if selectedRole.Name(event.State()) == "" {
		return true
	}

	hasRole, err := p.userHasRole(event, selectedRole.ServerRoleID)
	if err != nil {
		event.Except(err)
		return true
	}

	if hasRole {
		p.removeRole(event, event.MessageReactionAdd.ChannelID, selectedRole.ServerRoleID)
	} else {
		if selectedRole.CategoryID == 0 {
			p.assignRole(event, event.MessageReactionAdd.ChannelID, selectedRole.ServerRoleID)
		} else {
			categories, err := p.getCategoryByChannel(event.ChannelID)
			if err != nil {
				event.ExceptSilent(err)
				return true
			}

			var selectedCategory *Category
		CatLoop:
			for _, cat := range categories {
				for _, crole := range cat.Roles {
					if crole.ServerRoleID == selectedRole.ServerRoleID {
						selectedCategory = cat
						break CatLoop
					}
				}
			}

			if selectedCategory == nil {
				return true
			}

			member, err := event.State().Member(event.GuildID, event.UserID)
			if err != nil {
				event.ExceptSilent(err)
				return false
			}

			if p.isOverRoleLimit(member, selectedCategory) {
				msgs, err := event.Send(event.ChannelID, "roles.role.at-category-limit", "userMention", member.Mention())
				if err != nil {
					return false
				}

				go p.deleteWithDelay(event, msgs[0].ID)
				return true
			}
			p.assignRole(event, event.MessageReactionAdd.ChannelID, selectedRole.ServerRoleID)
		}
	}

	return true
}

func (p *Plugin) handleUserRoleRequest(event *events.Event) bool {
	if event.MessageCreate == nil || event.MessageCreate.Author == nil || event.MessageCreate.Author.Bot || !p.isInRoleChannel(event) {
		return false
	}

	if len(event.MessageCreate.Content) < 2 {
		go p.deleteWithDelay(event, event.MessageID)
		return true
	}
	plusMinus := event.MessageCreate.Content[0:1]

	if event.HasOr(permissions.DiscordAdministrator, permissions.DiscordManageRoles) && plusMinus != PLUS && plusMinus != MINUS {
		return false
	}

	go p.deleteWithDelay(event, event.MessageID)

	if event.Command() {
		return true
	}

	// check plus or minus
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

	// Try case sensitive matching first, then try case insensitive
	requests, err := p.parseRoleRequestMessage(event, strings.TrimSpace(event.MessageCreate.Content), allRoles, true)
	if err != nil {
		requests, err = p.parseRoleRequestMessage(event, strings.TrimSpace(event.MessageCreate.Content), allRoles, false)
	}
	if err != nil {
		event.ExceptSilent(err)
		msgs, err := event.Send(event.ChannelID, err.Error())
		if err == nil {
			go p.deleteWithDelay(event, msgs[0].ID)
		}
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
						err = p.assignRole(event, event.ChannelID, role.ServerRoleID)
					} else {
						err = p.removeRole(event, event.ChannelID, role.ServerRoleID)
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

						err = p.assignRole(event, event.ChannelID, role.ServerRoleID)
					} else {
						err = p.removeRole(event, event.ChannelID, role.ServerRoleID)
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

func (p *Plugin) parseRoleRequestMessage(event *events.Event, userMsg string, roles []*Role, caseSensitive bool) (map[*Role]rune, error) {
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
				if len(roleName) != 0 && hasPrefix(remainingMsg, roleName, caseSensitive) {
					ignoreChars = len(roleName)
					requestsMap[role] = str
					foundRole = true
					skipSpace = true
					break RoleLoop
				}

				for _, alias := range role.Aliases {
					if len(alias) != 0 && hasPrefix(remainingMsg, alias, caseSensitive) {
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

func hasPrefix(s string, prefix string, caseSensitive bool) bool {
	if caseSensitive {
		return strings.HasPrefix(s, prefix)
	}
	return strings.HasPrefix(strings.ToLower(s), strings.ToLower(prefix))
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

func (p *Plugin) assignRole(event *events.Event, channelID string, serverRoleID string) error {
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
			msgs, err := event.Send(channelID, "roles.role.already-assigned", "userMention", member.Mention(), "serverRoleName", role.Name)
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

	msgs, err := event.Send(channelID, "roles.role.assigned", "userMention", member.Mention(), "serverRoleName", role.Name)
	if err != nil {
		return err
	}
	go p.deleteWithDelay(event, msgs[0].ID)

	return nil
}

func (p *Plugin) removeRole(event *events.Event, channelID string, serverRoleID string) error {
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
		msgs, err := event.Send(channelID, "roles.role.not-assigned", "userMention", member.Mention(), "serverRoleName", role.Name)
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

	msgs, err := event.Send(channelID, "roles.role.removed-role", "userMention", member.Mention(), "serverRoleName", role.Name)
	if err != nil {
		return err
	}

	go p.deleteWithDelay(event, msgs[0].ID)

	return nil
}

func (p *Plugin) userHasRole(event *events.Event, serverRoleID string) (bool, error) {
	// confirm the user has the role
	member, err := event.State().Member(event.GuildID, event.UserID)
	if err != nil {
		return false, err
	}

	for _, userRole := range member.Roles {
		if userRole == serverRoleID {
			return true, nil
		}
	}

	return false, nil
}

func (p *Plugin) isInRoleChannel(event *events.Event) bool {
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
	return inRoleChannel
}
