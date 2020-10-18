package roles

import (
	"strings"
	"time"

	"gitlab.com/Cacophony/go-kit/config"
	"go.uber.org/zap"
)

const (
	INTERVAL = 30 * time.Minute
)

func (p *Plugin) startRoleChannelCacheLoop() {
	err := p.cacheAllRoleChannels()
	if err != nil {
		return
	}

	go func() {
		var err error
		for {
			time.Sleep(INTERVAL)

			err = p.cacheAllRoleChannels()
			if err != nil {
				p.logger.Error("failed to cache role channels", zap.Error(err))
				continue
			}

			p.logger.Debug("cached role channels")
		}
	}()
}

func (p *Plugin) cacheAllRoleChannels() error {

	var categories []Category
	err := p.db.Find(&categories).Error
	if err != nil && !strings.Contains(err.Error(), "record not found") {
		return err
	}

	channels := make(map[string][]string)

	for _, cat := range categories {
		if cat.ChannelID != "" {
			channels[cat.GuildID] = append(channels[cat.GuildID], cat.ChannelID)
		}
	}

	type Result struct {
		GuildID string
	}
	var results []Result
	err = p.db.
		Table((&Role{}).TableName()).
		Select("distinct guild_id").
		Find(&results).
		Error
	if err != nil {
		p.logger.Error("no guilds found with roles during cache", zap.Error(err))
		return err
	}

	for _, result := range results {
		channelID, err := config.GuildGetString(p.db, result.GuildID, guildRoleChannelKey)
		if err != nil || channelID == "" {
			continue
		}
		channels[result.GuildID] = append(channels[result.GuildID], channelID)
	}

	p.guildRoleChannelsLock.Lock()
	p.guildRoleChannels = channels
	p.guildRoleChannelsLock.Unlock()
	return nil
}

func (p *Plugin) cacheGuildsRoleChannels(guildID string) error {

	var categories []Category
	err := p.db.Where("guild_id = ?", guildID).Find(&categories).Error
	if err != nil && !strings.Contains(err.Error(), "record not found") {
		return err
	}

	channels := make([]string, 0)

	for _, cat := range categories {
		if cat.ChannelID != "" {
			channels = append(channels, cat.ChannelID)
		}
	}

	defaultRoleChannel, err := config.GuildGetString(p.db, guildID, guildRoleChannelKey)
	if err == nil && defaultRoleChannel != "" {
		channels = append(channels, defaultRoleChannel)
	}

	p.guildRoleChannelsLock.Lock()
	p.guildRoleChannels[guildID] = channels
	p.guildRoleChannelsLock.Unlock()
	return nil
}

func (p *Plugin) getCachedRoleChannels(guildID string) []string {
	if p.guildRoleChannels == nil {
		return nil
	}
	p.guildRoleChannelsLock.RLock()
	channels := p.guildRoleChannels[guildID]
	p.guildRoleChannelsLock.RUnlock()
	return channels
}
