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

func (h *Handler) Handle(event *events.Event) (process bool) {
	var err error

	process = false

	if event.GuildID == "" {
		return
	}

	if event.Type == events.CacophonyAutomodWait {
		h.handleWaitEvent(event)
		return true
	}

	h.rulesLock.RLock()
	rules, ok := h.rules[event.GuildID]
	h.rulesLock.RUnlock()
	if !ok {
		return
	}

	var triggerMatched bool
	for _, rule := range rules {
		rule := rule

		env := &models.Env{}
		h.addBaseToEnv(env, event)
		env.Rule = &rule

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

		doProceed := h.executeFilters(env, &rule)
		if !doProceed {
			continue
		}

		log, runError := h.executeActions(env, &rule)

		if log {
			err = h.logRun(env, rule, runError)
			if err != nil &&
				err != state.ErrBotForGuildStateNotFound {
				event.ExceptSilent(err)
			}
		}

		if rule.Stop {
			process = true
		}
	}

	return
}

func (h *Handler) handleWaitEvent(event *events.Event) {
	env := &models.Env{}
	err := env.Unmarshal(event.AutomodWait.EnvData)
	if err != nil {
		env.Event.ExceptSilent(err)
		return
	}
	h.addBaseToEnv(env, event)

	doProceed := h.executeFilters(env, env.Rule)
	if !doProceed {
		return
	}

	log, runError := h.executeActions(env, env.Rule)

	if log {
		err = h.logRun(env, *env.Rule, runError)
		if err != nil &&
			err != state.ErrBotForGuildStateNotFound {
			event.ExceptSilent(err)
		}
	}
}

func (h *Handler) addBaseToEnv(env *models.Env, event *events.Event) {
	env.Event = event
	env.State = h.state
	env.Redis = h.redis
	env.Handler = h // TODO: remove this field
	env.Tokens = h.tokens
}

func (h *Handler) executeFilters(env *models.Env, rule *models.Rule) bool {
	for _, filter := range list.FiltersList {
		for _, ruleFilter := range rule.Filters {
			if filter.Name() != ruleFilter.Name {
				continue
			}

			item, err := filter.NewItem(env, ruleFilter.Values)
			if err != nil {
				env.Event.ExceptSilent(err)
				return false
			}

			if ruleFilter.Not {
				if item.Match(env) {
					return false
				}
			} else {
				if !item.Match(env) {
					return false
				}
			}
		}
	}

	return true
}

func (h *Handler) executeActions(env *models.Env, rule *models.Rule) (bool, error) {
	var runError error

	for _, ruleAction := range rule.Actions {
		for _, action := range list.ActionsList {
			if action.Name() != ruleAction.Name {
				continue
			}

			item, err := action.NewItem(env, ruleAction.Values)
			if err != nil {
				env.Event.ExceptSilent(err)
				return true, err
			}

			stop, err := item.Do(env)
			if err != nil {
				runError = err
			}
			if stop {
				// do not log waits without error, because they will be logged later
				if action.Name() == "wait" && err == nil {
					return false, runError
				}
				return true, runError
			}
		}
	}

	return true, runError
}
