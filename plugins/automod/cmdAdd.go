package automod

import (
	"encoding/json"
	"strings"

	"gitlab.com/Cacophony/Processor/plugins/automod/handler"
	"gitlab.com/Cacophony/Processor/plugins/automod/list"
	"gitlab.com/Cacophony/Processor/plugins/automod/models"
	"gitlab.com/Cacophony/go-kit/config"
	"gitlab.com/Cacophony/go-kit/discord"
	"gitlab.com/Cacophony/go-kit/events"
)

const confirmUpdateRuleQuestionnaireKey = "cacophony:processor:automod:confirm-update-rule"

func (p *Plugin) cmdAdd(event *events.Event) {
	if len(event.Fields()) < 3 {
		event.Respond("automod.add.too-few")
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
			event.Respond("automod.add.invalid-trigger-value", "error", err)
			return
		}
		newRule.TriggerName = triggerName
		newRule.TriggerValues = triggerArgs
		break
	}
	if newRule.TriggerName == "" {
		event.Respond("automod.add.invalid-trigger-name")
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
				event.Respond("automod.add.invalid-filter-value", "error", err)
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
		event.Respond("automod.add.invalid-filter-name")
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
				event.Respond("automod.add.invalid-action-value", "error", err)
				return
			}
			newRule.Actions = append(newRule.Actions, models.RuleAction{
				Name:   actionName,
				Values: actionArgs,
			})
			break
		}
	}
	if len(newRule.Actions) == 0 {
		event.Respond("automod.add.invalid-action-name")
		return
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

	var existingRulesWithName []models.Rule
	err = p.db.Model(models.Rule{}).Where(models.Rule{
		Name:    newRule.Name,
		GuildID: event.GuildID,
	}).Find(&existingRulesWithName).Error
	if err != nil {
		event.Except(err)
		return
	}
	if len(existingRulesWithName) > 0 {
		if existingRulesWithName[0].Managed {
			event.Respond("automod.add.managed-duplicate")
			return
		}

		messages, err := event.Respond("automod.add.confirm-update-duplicate")
		if err != nil {
			event.Except(err)
			return
		}

		ruleData, err := json.Marshal(&newRule)
		if err != nil {
			event.Except(err)
			return
		}

		err = event.Questionnaire().Register(
			confirmUpdateRuleQuestionnaireKey,
			events.QuestionnaireFilter{
				GuildID:   event.GuildID,
				ChannelID: event.ChannelID,
				UserID:    event.UserID,
				Type:      events.MessageReactionAddType,
			},
			map[string]interface{}{
				"messageID": messages[0].ID,
				"rule":      string(ruleData),
			},
		)
		if err != nil {
			event.Except(err)
			return
		}

		err = discord.React(
			event.Redis(), event.Discord(), messages[0].ChannelID, messages[0].ID, false, "✅",
		)
		if err != nil {
			return
		}
		discord.React(
			event.Redis(), event.Discord(), messages[0].ChannelID, messages[0].ID, false, "❌",
		)

		return
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

	if len(fields) < 1 {
		return "", nil, false, fieldsInput
	}

	for _, filter := range list.FiltersList {
		if filter.Name() != fields[0] {
			continue
		}

		name = filter.Name()
		if len(fields) > filter.Args() {
			args = fields[1 : filter.Args()+1]
		}
		if len(fields) > filter.Args() {
			fields = fields[filter.Args()+1:]
		}
		return
	}

	return "", nil, false, fieldsInput
}

func extractAction(fieldsInput []string) (name string, args []string, fields []string) {
	if len(fieldsInput) < 1 {
		return "", nil, fieldsInput
	}

	for _, action := range list.ActionsList {
		if action.Name() != fieldsInput[0] {
			continue
		}

		name = action.Name()
		if len(fieldsInput) > action.Args() {
			args = fieldsInput[1 : action.Args()+1]
		}
		if len(fieldsInput) > action.Args() {
			fields = fieldsInput[action.Args()+1:]
		}
		return
	}

	return "", nil, fieldsInput
}
