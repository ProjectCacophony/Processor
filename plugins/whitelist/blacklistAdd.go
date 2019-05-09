package whitelist

import (
	"strings"

	"go.uber.org/zap"

	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) blacklistAdd(event *events.Event) {
	if len(event.Fields()) < 3 {
		event.Respond("whitelist.blacklist-add.too-few-args")
		return
	}

	guild, err := p.extractGuild(event.Discord(), event.Fields()[2])
	if err != nil {
		event.Respond("whitelist.blacklist-add.invalid-invite")
		return
	}

	blacklistEntry, err := blacklistFind(p.db, "guild_id = ?", guild.ID)
	if err != nil {
		if !strings.Contains(err.Error(), "record not found") {
			event.Except(err)
			return
		}
	}

	if blacklistEntry != nil {
		event.Respond("whitelist.blacklist-add.already-blacklisted")
		return
	}

	err = whitelistRemove(p.db, guild.ID)
	if err != nil {
		event.Except(err)
		return
	}

	err = blacklistAdd(p.db, event.UserID, guild.ID)
	if err != nil {
		event.Except(err)
		return
	}

	p.logger.Info(
		"blacklisted server",
		zap.String("user_id", event.UserID),
		zap.String("guild_id", guild.ID),
	)

	err = cacheWhitelistAndBlacklist(p.db, p.redis)
	if err != nil {
		event.Except(err)
		return
	}

	_, err = event.Respond("blacklist.add.success", "name", guild.Name)
	event.Except(err)
}
