package autoroles

import (
	"fmt"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/Processor/plugins/automod/filters"
	"gitlab.com/Cacophony/Processor/plugins/automod/models"
	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) createAutoRole(event *events.Event) {
	if len(event.Fields()) < 3 {
		event.Respond("common.invalid-params")
		return
	}

	var duration int
	var err error

	// check for duration
	if len(event.Fields()) == 4 {
		duration, err = strconv.Atoi(event.Fields()[3])
		if err != nil {
			event.Respond("common.invalid-params")
			return
		}
	}

	serverRoleID := event.Fields()[2]
	if serverRoleID == "" {
		event.Respond("roles.role.role-not-found-on-server")
		return
	}

	serverRole, err := p.getServerRoleByNameOrID(serverRoleID, event.GuildID)
	if err != nil {
		event.Except(err)
		return
	}

	rule := newAutoRoleRule(event, time.Duration(duration)*time.Second, serverRole)

	err = models.CreateRule(event.DB(), rule)
	if err != nil {
		event.Except(err)
		return
	}

	autoRole := &AutoRole{
		GuildID:      rule.GuildID,
		ServerRoleID: serverRole.ID,
		RuleID:       rule.ID,
	}

	err = p.db.Save(autoRole).Error
	if err != nil {
		models.DeleteRule(event.DB(), rule)
		event.Except(err)
		return
	}

	event.Respond("autorole.creation",
		"roleName", serverRole.Name,
	)
}

func (p *Plugin) deleteAutoRole(event *events.Event) {
	if len(event.Fields()) < 3 {
		event.Respond("common.invalid-params")
		return
	}

	serverRole, err := p.getServerRoleByNameOrID(event.Fields()[2], event.GuildID)
	if err != nil {
		event.Respond("roles.role.role-not-found-on-server")
		return
	}

	autoRole, err := p.getAutoRolesByServerRoleID(serverRole.ID, event.GuildID)
	if err != nil {
		event.Respond("autorole.not-found")
		return
	}

	err = p.db.
		Delete(autoRole.Rule.Actions, "rule_id = ?", autoRole.Rule.ID).
		Delete(autoRole.Rule).
		Delete(autoRole).Error
	if err != nil {
		event.Except(err)
		return
	}

	event.Respond("autorole.deleted",
		"roleName", serverRole.Name,
	)
}

func (p *Plugin) applyAutoRole(event *events.Event) {
	if len(event.Fields()) < 2 {
		event.Respond("common.invalid-params")
		return
	}

	autoRoles, err := p.getAutoRoles(event.GuildID)
	if err != nil {
		event.Except(err)
		return
	}

	userIDs, err := event.State().GuildMembers(event.GuildID)
	if err != nil {
		event.Respond("autorole.guild-members-not-found")
		return
	}

	for _, aRole := range autoRoles {

		var delay time.Duration
		for _, action := range aRole.Rule.Actions {
			if action.Name == (filters.Wait{}).Name() {
				delay, err = time.ParseDuration(action.Values[0])
				if err != nil {
					event.Except(err)
					return
				}
			}
		}

		for _, userID := range userIDs {
			member, err := event.State().Member(event.GuildID, userID)
			if err != nil {
				continue
			}

			// check if user has been in the server for less time than the delay
			if delay > 0 {
				join, err := member.JoinedAt.Parse()
				if err != nil {
					continue
				}

				if time.Since(join) < delay {
					continue
				}
			}

			hasRole := false
			for _, userRoleID := range member.Roles {
				if userRoleID == aRole.ServerRoleID {
					hasRole = true
					break
				}
			}

			if !hasRole {
				err := event.Discord().Client.GuildMemberRoleAdd(event.GuildID, userID, aRole.ServerRoleID)
				if err != nil {
					event.ExceptSilent(err)
				}
			}
		}
	}

	event.Respond("autorole.applied")
}

func (p *Plugin) listAutoRoles(event *events.Event) {
	autoRoles, err := p.getAutoRoles(event.GuildID)
	if err != nil || len(autoRoles) == 0 {
		event.Respond("autorole.none-found")
		return
	}

	msg := "**__Auto Roles__**"
	for _, aRole := range autoRoles {

		serverRole, err := p.getServerRoleByNameOrID(aRole.ServerRoleID, event.GuildID)
		if err != nil {
			msg += "*Unknown*"
			continue
		}

		var delay string
		for _, action := range aRole.Rule.Actions {
			if action.Name == (filters.Wait{}).Name() {
				delay = action.Values[0]
			}
		}

		if delay == "" {
			msg += fmt.Sprintf("\n**%s** (%s)", serverRole.Name, serverRole.ID)
			continue
		}

		msg += fmt.Sprintf("\n**%s** - %s delay (%s)", serverRole.Name, delay, serverRole.ID)
	}

	event.Respond(msg)
}

func newAutoRoleRule(
	event *events.Event,
	delay time.Duration,
	role *discordgo.Role,
) *models.Rule {
	rule := &models.Rule{
		GuildID: event.GuildID,
		Name:    "Auto Role",
		Actions: []models.RuleAction{
			{
				Name: "apply_role",
				Values: []string{
					role.ID,
				},
			},
		},
		TriggerName: "when_join",
		Silent:      true,
		Managed:     true,
	}

	if delay.Seconds() > 0 {
		rule.Actions = append(
			[]models.RuleAction{{
				Name: "wait",
				Values: []string{
					delay.String(),
				},
			}},
			rule.Actions...,
		)
	}

	return rule
}
