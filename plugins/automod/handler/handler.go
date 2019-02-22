package handler

import (
	"sync"

	"github.com/go-redis/redis"

	"github.com/jinzhu/gorm"
	"gitlab.com/Cacophony/Processor/plugins/automod/list"
	"gitlab.com/Cacophony/Processor/plugins/automod/models"
	"gitlab.com/Cacophony/go-kit/events"
	"go.uber.org/zap"
)

type Handler struct {
	logger *zap.Logger
	db     *gorm.DB
	redis  *redis.Client

	rules     map[string][]models.Rule
	rulesLock sync.RWMutex
}

func NewHandler(logger *zap.Logger, db *gorm.DB, redis *redis.Client) (*Handler, error) {
	handler := &Handler{
		logger: logger,
		db:     db,
		redis:  redis,
	}

	err := handler.startRulesCaching()

	return handler, err
}

func (h *Handler) Handle(event *events.Event) (process bool) {
	env := &models.Env{
		Event: event,
		State: event.State(),
		Redis: h.redis,
	}

	if event.Type == events.MessageCreateType {
		if event.MessageCreate.Author.Bot {
			return false
		}

		env.GuildID = event.MessageCreate.GuildID
		env.UserID = event.MessageCreate.Author.ID
	}

	h.rulesLock.RLock()
	rules, ok := h.rules[env.GuildID]
	h.rulesLock.RUnlock()
	if !ok {
		return true
	}

	var triggerMatched bool
	var filtersMatched bool
	for _, rule := range rules {
		triggerMatched = false

		for _, trigger := range list.TriggerList {
			if trigger.Name() != rule.TriggerName {
				continue
			}

			item, err := trigger.NewItem(env, rule.TriggerValues)
			if err != nil {
				event.ExceptSilent(err)
				return false
			}

			if item.Match(env) {
				triggerMatched = true
			}
		}

		if !triggerMatched {
			continue
		}

		filtersMatched = true

		for _, filter := range list.FiltersList {
			for _, ruleFilter := range rule.Filters {
				if filter.Name() != ruleFilter.Name {
					continue
				}

				item, err := filter.NewItem(env, ruleFilter.Values)
				if err != nil {
					event.ExceptSilent(err)
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

		for _, action := range list.ActionsList {
			for _, ruleAction := range rule.Actions {
				if action.Name() != ruleAction.Name {
					continue
				}

				item, err := action.NewItem(env, ruleAction.Values)
				if err != nil {
					event.ExceptSilent(err)
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
