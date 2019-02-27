package whitelist

import (
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
)

func (p *Plugin) startWhitelistAndBlacklistCaching() error {
	err := cacheWhitelistAndBlacklist(p.db, p.redis)
	if err != nil {
		return err
	}

	go func() {
		var err error
		for {
			time.Sleep(interval)

			err = cacheWhitelistAndBlacklist(p.db, p.redis)
			if err != nil {
				p.logger.Error("failed to cache whitelist and blacklist",
					zap.Error(err),
				)
			}

			p.logger.Debug("cached whitelist and blacklist")
		}
	}()

	return nil
}

const (
	whitelistKey = "cacophony.whitelist.whitelist"
	blacklistKey = "cacophony.whitelist.blacklist"

	expiration = time.Hour * 24 * 7 // one week
	interval   = time.Hour
)

func cacheWhitelistAndBlacklist(db *gorm.DB, redis *redis.Client) error {
	err := cacheWhitelist(db, redis)
	if err != nil {
		return err
	}

	err = cacheBlacklist(db, redis)
	return err
}

func cacheWhitelist(db *gorm.DB, redis *redis.Client) error {
	servers, err := whitelistAll(db)
	if err != nil {
		return err
	}

	guildIDs := make([]string, len(servers))
	for i, server := range servers {
		guildIDs[i] = server.GuildID
	}

	return redis.Set(whitelistKey, strings.Join(guildIDs, ";"), expiration).Err()
}

func cacheBlacklist(db *gorm.DB, redis *redis.Client) error {
	servers, err := blacklistAll(db)
	if err != nil {
		return err
	}

	guildIDs := make([]string, len(servers))
	for i, server := range servers {
		guildIDs[i] = server.GuildID
	}

	return redis.Set(blacklistKey, strings.Join(guildIDs, ";"), expiration).Err()
}
