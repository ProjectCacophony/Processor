package automod

import (
	"strings"

	"gitlab.com/Cacophony/Processor/plugins/automod/models"
	"gitlab.com/Cacophony/go-kit/events"
)

func cmdAdd(event *events.Event) {
	if len(event.Fields()) < 8 {
		event.Respond("automod.add.too-few") // nolint: errcheck
		return
	}

	// TODO: support adding multiple triggers

	var err error
	var newRule Rule
	newRule.GuildID = event.MessageCreate.GuildID
	newRule.Name = event.Fields()[2]

	env := &models.Env{
		State:   event.State(),
		GuildID: event.MessageCreate.GuildID,
		UserID:  event.MessageCreate.Author.ID,
	}

	for _, trigger := range triggerList {
		if trigger.Name() != event.Fields()[3] {
			continue
		}

		newRule.Trigger = trigger.Name()
	}
	if newRule.Trigger == "" {
		event.Respond("automod.add.invalid-trigger-name") // nolint: errcheck
		return
	}

	var newFilter RuleFilter
	for _, filter := range filtersList {
		if filter.Name() != event.Fields()[4] {
			continue
		}

		_, err = filter.NewItem(env, event.Fields()[5])
		if err != nil {
			event.Respond("automod.add.invalid-filter-value", "error", err) // nolint: errcheck
			return
		}

		newFilter.Name = filter.Name()
		newFilter.Value = event.Fields()[5]
	}
	if newFilter.Name == "" {
		event.Respond("automod.add.invalid-filter-name") // nolint: errcheck
		return
	}
	newRule.Filters = []RuleFilter{newFilter}

	var newAction RuleAction
	for _, action := range actionsList {
		if action.Name() != event.Fields()[6] {
			continue
		}

		_, err = action.NewItem(env, event.Fields()[7])
		if err != nil {
			event.Respond("automod.add.invalid-action-value", "error", err) // nolint: errcheck
			return
		}

		newAction.Name = action.Name()
		newAction.Value = event.Fields()[7]
	}
	if newAction.Name == "" {
		event.Respond("automod.add.invalid-action-name") // nolint: errcheck
		return
	}
	newRule.Actions = []RuleAction{newAction}

	if event.Fields()[len(event.Fields())-1] == "continue" {
		newRule.Process = true
	}

	err = event.DB().Save(&newRule).Error
	if err != nil {
		if strings.Contains(err.Error(), "idx_automod_rules_guildid_name") {
			event.Respond("automod.add.name-in-use") // nolint: errcheck
			return
		}
		event.Except(err)
		return
	}

	_, err = event.Respond("automote.add.success")
	event.Except(err)
}
