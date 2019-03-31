package automod

import (
	"strings"

	"gitlab.com/Cacophony/Processor/plugins/automod/handler"

	"gitlab.com/Cacophony/Processor/plugins/automod/list"
	"gitlab.com/Cacophony/Processor/plugins/automod/models"
	"gitlab.com/Cacophony/go-kit/config"
	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) cmdAdd(event *events.Event) {
	if len(event.Fields()) < 3 {
		event.Respond("automod.add.too-few") // nolint: errcheck
		return
	}

	var newRule models.Rule
	newRule.GuildID = event.MessageCreate.GuildID
	newRule.Name = event.Fields()[2]
	fields := event.Fields()[3:]

	env := &models.Env{
		State:   p.state,
		GuildID: event.MessageCreate.GuildID,
		UserID:  []string{event.MessageCreate.Author.ID},
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
	var not bool
	for {
		filterName, filterArgs, not, fields = extractFilter(fields)
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
				Not:    not,
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

		switch field {
		case "stop":
			newRule.Stop = true
		case "silent":
			newRule.Silent = true
		}
	}

	var logChannelSet string
	_, err := config.GuildGetString(p.db, event.GuildID, handler.AutomodLogKey)
	if err != nil {
		if !strings.Contains(err.Error(), "record not found") {
			event.Except(err)
			return
		}

		logChannelSet = event.ChannelID
		err = config.GuildSetString(p.db, event.GuildID, handler.AutomodLogKey, logChannelSet)
		if err != nil {
			event.Except(err)
			return
		}
	}

	err = p.db.Save(&newRule).Error
	if err != nil {
		event.Except(err)
		return
	}

	_, err = event.Respond("automote.add.success",
		"logChannelID", logChannelSet,
	)
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

func extractFilter(fieldsInput []string) (name string, args []string, not bool, fields []string) {
	fields = fieldsInput

	if len(fields) >= 1 && fields[0] == "not" {
		not = true
		fields = fields[1:]
	}

	for i, field := range fields {
		for _, filter := range list.FiltersList {
			if filter.Name() != field {
				continue
			}

			name = filter.Name()
			if len(fields) > i+filter.Args() {
				args = fields[i+1 : i+filter.Args()+1]
			}
			if len(fields) > i+filter.Args() {
				fields = fields[i+filter.Args()+1:]
			}
			return
		}
	}

	return "", nil, false, fieldsInput
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
