package handler

import (
	"sync"

	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"gitlab.com/Cacophony/Processor/plugins/automod/list"
	"gitlab.com/Cacophony/Processor/plugins/automod/models"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/state"
	"go.uber.org/zap"
)

type Handler struct {
	logger *zap.Logger
	db     *gorm.DB
	redis  *redis.Client
	tokens map[string]string
	state  *state.State

	rules           map[string][]models.Rule
	rulesLock       sync.RWMutex
	logChannels     map[string]string
	logChannelsLock sync.RWMutex
}

func NewHandler(
	logger *zap.Logger,
	db *gorm.DB,
	redis *redis.Client,
	tokens map[string]string,
	state *state.State,
) (*Handler, error) {
	handler := &Handler{
		logger: logger,
		db:     db,
		redis:  redis,
		tokens: tokens,
		state:  state,
	}

	err := handler.startRulesCaching()
	if err != nil {
		return nil, err
	}

	err = handler.startLogChannelsCaching()
	return handler, err
}

// nolint: nakedret
func (h *Handler) Handle(event *events.Event) (process bool) {
	process = false

	env := &models.Env{
		Event:   event,
		State:   h.state,
		Redis:   h.redis,
		Handler: h,
		Tokens:  h.tokens,
	}

	if event.GuildID == "" {
		return
	}

	h.rulesLock.RLock()
	rules, ok := h.rules[event.GuildID]
	h.rulesLock.RUnlock()
	if !ok {
		return
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
				return
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
					return
				}

				if ruleFilter.Not {
					if item.Match(env) {
						filtersMatched = false
					}
				} else {
					if !item.Match(env) {
						filtersMatched = false
					}
				}
			}
		}

		if !filtersMatched {
			continue
		}

		if !rule.Silent {
			err := h.postLog(env, rule)
			if err != nil {
				event.ExceptSilent(err)
			}
		}

		for _, action := range list.ActionsList {
			for _, ruleAction := range rule.Actions {
				if action.Name() != ruleAction.Name {
					continue
				}

				item, err := action.NewItem(env, ruleAction.Values)
				if err != nil {
					event.ExceptSilent(err)
					return
				}

				item.Do(env)
			}
		}

		if rule.Stop {
			process = true
		}
	}

	return
}
