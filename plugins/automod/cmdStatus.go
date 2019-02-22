package automod

import (
	"strings"

	"gitlab.com/Cacophony/Processor/plugins/automod/models"
	"gitlab.com/Cacophony/go-kit/events"
)

func cmdStatus(event *events.Event) {
	var rules []models.Rule
	err := event.DB().
		Preload("Filters").
		Preload("Actions").
		Where("guild_id = ?", event.MessageCreate.GuildID).
		Find(&rules).Error
	if err != nil {
		event.Except(err)
		return
	}

	ruleTexts := make([]string, len(rules))
	for i, rule := range rules {
		ruleTexts[i] += addQuotesIfSpaces(rule.Name) + " "
		ruleTexts[i] += addQuotesIfSpaces(rule.TriggerName) + " "
		ruleTexts[i] += argsString(rule.TriggerValues)
		for _, filter := range rule.Filters {
			ruleTexts[i] += addQuotesIfSpaces(filter.Name) + " "
			ruleTexts[i] += argsString(filter.Values)
		}
		for _, action := range rule.Actions {
			ruleTexts[i] += addQuotesIfSpaces(action.Name) + " "
			ruleTexts[i] += argsString(action.Values)
		}
		if rule.Process {
			ruleTexts[i] += "continue "
		}
	}

	_, err = event.Respond("automod.status.response", "ruleTexts", ruleTexts)
	event.Except(err)
}

func addQuotesIfSpaces(input string) string {
	if strings.Contains(input, " ") {
		return "\"" + input + "\""
	}

	return input
}

func argsString(input []string) string {
	var result string
	for _, item := range input {
		result += addQuotesIfSpaces(item) + " "
	}
	return result
}
