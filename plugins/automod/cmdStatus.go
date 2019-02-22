package automod

import (
	"sort"
	"strings"

	"gitlab.com/Cacophony/Processor/plugins/automod/models"
	"gitlab.com/Cacophony/go-kit/events"
)

type sortRulesByName []models.Rule

// Len is part of sort.Interface
func (d sortRulesByName) Len() int {
	return len(d)
}

// Swap is part of sort.Interface
func (d sortRulesByName) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

// Less is part of sort.Interface
func (d sortRulesByName) Less(i, j int) bool {
	return strings.ToLower(d[i].Name) < strings.ToLower(d[j].Name)
}

func (p *Plugin) cmdStatus(event *events.Event) {
	var rules []models.Rule
	err := p.db.
		Preload("Filters").
		Preload("Actions").
		Where("guild_id = ?", event.MessageCreate.GuildID).
		Find(&rules).Error
	if err != nil {
		event.Except(err)
		return
	}

	sort.Sort(sortRulesByName(rules))

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
