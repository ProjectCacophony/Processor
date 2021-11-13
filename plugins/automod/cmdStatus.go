package automod

import (
	"sort"
	"strings"

	"gitlab.com/Cacophony/Processor/plugins/automod/handler"
	"gitlab.com/Cacophony/Processor/plugins/automod/models"
	"gitlab.com/Cacophony/go-kit/config"
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

type enhancedRule struct {
	models.Rule
	RuleText string
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

	var ruleText string
	enhancedRules := make([]enhancedRule, len(rules))
	for i, rule := range rules {
		ruleText = ""
		ruleText += addQuotesIfSpaces(rule.Name) + " "
		ruleText += addQuotesIfSpaces(rule.TriggerName) + " "
		ruleText += argsString(rule.TriggerValues)
		for _, filter := range rule.Filters {
			if filter.Not {
				ruleText += "not "
			}
			ruleText += addQuotesIfSpaces(filter.Name) + " "
			ruleText += argsString(filter.Values)
		}
		for _, action := range rule.Actions {
			ruleText += addQuotesIfSpaces(action.Name) + " "
			ruleText += argsString(action.Values)
		}
		if rule.Stop {
			ruleText += "stop "
		}
		if rule.Silent {
			ruleText += "silent "
		}

		enhancedRules[i].Rule = rule
		enhancedRules[i].RuleText = ruleText
	}

	logChannel, err := config.GuildGetString(p.db, event.GuildID, handler.AutomodLogKey)
	if err != nil && !strings.Contains(err.Error(), "record not found") {
		event.Except(err)
		return
	}

	_, err = event.Respond(
		"automod.status.response",
		"rules", enhancedRules,
		"logChannelID", logChannel,
	)
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
