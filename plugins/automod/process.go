package automod

import (
	"gitlab.com/Cacophony/Processor/plugins/automod/models"
	"gitlab.com/Cacophony/go-kit/events"
)

func handle(event *events.Event) bool {
	env := &models.Env{
		Event: event,
		State: event.State(),
	}

	var guildID string

	if event.Type == events.MessageCreateType {
		if event.MessageCreate.Author.Bot {
			return false
		}

		guildID = event.MessageCreate.GuildID

		env.GuildID = event.MessageCreate.GuildID
		env.UserID = event.MessageCreate.Author.ID
	}

	// TODO: cache rules
	var rules []Rule
	err := event.DB().
		Preload("Filters").
		Preload("Actions").
		Where("guild_id = ?", guildID).
		Find(&rules).Error
	if err != nil {
		event.Except(err) // TODO: handle errors silently
		return false
	}

	var triggerMatched bool
	var filtersMatched bool
	for _, rule := range rules {
		triggerMatched = false

		for _, trigger := range triggerList {
			if trigger.Name() != rule.Trigger {
				continue
			}

			item := trigger.NewItem(env)

			if item.Match(env) {
				triggerMatched = true
			}
		}

		if !triggerMatched {
			continue
		}

		filtersMatched = true

		for _, filter := range filtersList {
			for _, ruleFilter := range rule.Filters {
				if filter.Name() != ruleFilter.Name {
					continue
				}

				item, err := filter.NewItem(env, ruleFilter.Value)
				if err != nil {
					event.Except(err) // TODO: handle errors silently
					return false
				}

				if !item.Match(env) {
					filtersMatched = false
				}
			}
		}

		if !filtersMatched {
			continue
		}

		for _, action := range actionsList {
			for _, ruleAction := range rule.Actions {
				if action.Name() != ruleAction.Name {
					continue
				}

				item, err := action.NewItem(env, ruleAction.Value)
				if err != nil {
					event.Except(err) // TODO: handle errors silently
					return false
				}

				item.Do(env)
			}
		}

		if !rule.Process {
			return false
		}
	}

	return true
}
