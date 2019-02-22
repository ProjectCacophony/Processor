package automod

import (
	"strings"

	"gitlab.com/Cacophony/Processor/plugins/automod/list"
	"gitlab.com/Cacophony/Processor/plugins/automod/models"
	"gitlab.com/Cacophony/go-kit/events"
)

func cmdAdd(event *events.Event) {
	if len(event.Fields()) < 3 {
		event.Respond("automod.add.too-few") // nolint: errcheck
		return
	}

	var newRule models.Rule
	newRule.GuildID = event.MessageCreate.GuildID
	newRule.Name = event.Fields()[2]
	fields := event.Fields()[3:]

	env := &models.Env{
		State:   event.State(),
		GuildID: event.MessageCreate.GuildID,
		UserID:  event.MessageCreate.Author.ID,
	}

	triggerName, triggerArgs, fields := extractTrigger(fields)
	for _, trigger := range list.TriggerList {
		if trigger.Name() != triggerName {
			continue
		}

		_, err := trigger.NewItem(env, triggerArgs)
		if err != nil {
			event.Respond("automod.add.invalid-trigger-value", "error", err) // nolint: errcheck
			return
		}
		newRule.TriggerName = triggerName
		newRule.TriggerValues = triggerArgs
		break
	}
	if newRule.TriggerName == "" {
		event.Respond("automod.add.invalid-trigger-name") // nolint: errcheck
		return
	}

	var filterName string
	var filterArgs []string
	for {
		filterName, filterArgs, fields = extractFilter(fields)
		if filterName == "" {
			break
		}
		for _, filter := range list.FiltersList {
			if filter.Name() != filterName {
				continue
			}

			_, err := filter.NewItem(env, filterArgs)
			if err != nil {
				event.Respond("automod.add.invalid-filter-value", "error", err) // nolint: errcheck
				return
			}
			newRule.Filters = append(newRule.Filters, models.RuleFilter{
				Name:   filterName,
				Values: filterArgs,
			})
			break
		}
	}
	if len(newRule.Filters) == 0 {
		event.Respond("automod.add.invalid-filter-name") // nolint: errcheck
		return
	}

	var actionName string
	var actionArgs []string
	for {
		actionName, actionArgs, fields = extractAction(fields)
		if actionName == "" {
			break
		}
		for _, action := range list.ActionsList {
			if action.Name() != actionName {
				continue
			}

			_, err := action.NewItem(env, actionArgs)
			if err != nil {
				event.Respond("automod.add.invalid-action-value", "error", err) // nolint: errcheck
				return
			}
			newRule.Actions = append(newRule.Actions, models.RuleAction{
				Name:   actionName,
				Values: actionArgs,
			})
			break
		}
		if len(newRule.Actions) == 0 {
			event.Respond("automod.add.invalid-action-name") // nolint: errcheck
			return
		}
	}

	for _, field := range fields {
		if field != "continue" {
			continue
		}

		newRule.Process = true
		break
	}

	err := event.DB().Save(&newRule).Error
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

func extractTrigger(fieldsInput []string) (name string, args []string, fields []string) {
	for i, field := range fieldsInput {
		for _, trigger := range list.TriggerList {
			if trigger.Name() != field {
				continue
			}

			name = trigger.Name()
			if len(fieldsInput) > i+trigger.Args() {
				args = fieldsInput[i+1 : i+trigger.Args()+1]
			}
			if len(fieldsInput) > i+trigger.Args() {
				fields = fieldsInput[i+trigger.Args()+1:]
			}
			return
		}
	}

	return "", nil, fieldsInput
}

func extractFilter(fieldsInput []string) (name string, args []string, fields []string) {
	for i, field := range fieldsInput {
		for _, filter := range list.FiltersList {
			if filter.Name() != field {
				continue
			}

			name = filter.Name()
			if len(fieldsInput) > i+filter.Args() {
				args = fieldsInput[i+1 : i+filter.Args()+1]
			}
			if len(fieldsInput) > i+filter.Args() {
				fields = fieldsInput[i+filter.Args()+1:]
			}
			return
		}
	}

	return "", nil, fieldsInput
}

func extractAction(fieldsInput []string) (name string, args []string, fields []string) {
	for i, field := range fieldsInput {
		for _, action := range list.ActionsList {
			if action.Name() != field {
				continue
			}

			name = action.Name()
			if len(fieldsInput) > i+action.Args() {
				args = fieldsInput[i+1 : i+action.Args()+1]
			}
			if len(fieldsInput) > i+action.Args() {
				fields = fieldsInput[i+action.Args()+1:]
			}
			return
		}
	}

	return "", nil, fieldsInput
}
