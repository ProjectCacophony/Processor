package roles

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
	if len(event.Fields()) < 4 {
		event.Respond("common.invalid-params")
		return
	}

	var duration int
	var err error

	// check for duration
	if len(event.Fields()) == 5 {
		duration, err = strconv.Atoi(event.Fields()[4])
		if err != nil {
			event.Respond("common.invalid-params")
			return
		}
	}

	serverRoleID := event.Fields()[3]
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

	event.Respond("roles.autorole.creation",
		"roleName", serverRole.Name,
	)
}

func (p *Plugin) deleteAutoRole(event *events.Event) {
	if len(event.Fields()) < 4 {
		event.Respond("common.invalid-params")
		return
	}

	serverRole, err := p.getServerRoleByNameOrID(event.Fields()[3], event.GuildID)
	if err != nil {
		event.Respond("roles.role.role-not-found-on-server")
		return
	}

	autoRole, err := p.getAutoRolesByServerRoleID(serverRole.ID, event.GuildID)
	if err != nil {
		event.Respond("roles.autorole.not-found")
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

	event.Respond("roles.autorole.deleted",
		"roleName", serverRole.Name,
	)
}

func (p *Plugin) listAutoRoles(event *events.Event)  {
	autoRoles, err := p.getAutoRoles(event.GuildID)
	if err != nil || len(autoRoles) == 0 {
		event.Respond("roles.autorole.none-found")
		return
	}

	msg := fmt.Sprintf("**__Auto Roles__**")
	for _, aRole := range autoRoles {
		
		serverRole, err := p.getServerRoleByNameOrID(aRole.ServerRoleID, event.GuildID)
		if err != nil {
			msg += fmt.Sprintf("*Unknown*")
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
		Name:    fmt.Sprintf("Auto Role"),
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
