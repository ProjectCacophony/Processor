package handler

import (
	"time"

	"gitlab.com/Cacophony/Processor/plugins/automod/models"
	"go.uber.org/zap"
)

var cacheInterval = time.Minute

func (h *Handler) startRulesCaching() error {
	err := h.cacheRules()
	if err != nil {
		return err
	}

	go func() {
		var err error
		for {
			time.Sleep(cacheInterval)

			err = h.cacheRules()
			if err != nil {
				h.logger.Error("failed to cache rules", zap.Error(err))
			}

			h.logger.Debug("cached rules")
		}
	}()

	return nil
}

func (h *Handler) cacheRules() error {
	var rules []models.Rule
	err := h.db.
		Preload("Filters").
		Preload("Actions").
		Find(&rules).Error
	if err != nil {
		return err
	}

	rulesMap := make(map[string][]models.Rule)

	for _, rule := range rules {
		rulesMap[rule.GuildID] = append(rulesMap[rule.GuildID],
			rule,
		)
	}

	h.rulesLock.Lock()
	h.rules = rulesMap
	h.rulesLock.Unlock()

	return nil
}

func (h *Handler) startLogChannelsCaching() error {
	err := h.cacheLogChannels()
	if err != nil {
		return err
	}

	go func() {
		var err error
		for {
			time.Sleep(cacheInterval)

			err = h.cacheLogChannels()
			if err != nil {
				h.logger.Error("failed to cache log channels", zap.Error(err))
			}

			h.logger.Debug("cached log channels")
		}
	}()

	return nil
}

func (h *Handler) cacheLogChannels() error {
	channelsMap, err := h.getLogChannelIDs()
	if err != nil {
		return err
	}

	h.logChannelsLock.Lock()
	h.logChannels = channelsMap
	h.logChannelsLock.Unlock()

	return nil
}
