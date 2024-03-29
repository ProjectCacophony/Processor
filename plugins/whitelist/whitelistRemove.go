package whitelist

import (
	"strings"

	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/permissions"
	"go.uber.org/zap"
)

func (p *Plugin) whitelistRemove(event *events.Event) {
	if len(event.Fields()) < 3 {
		event.Respond("whitelist.remove.too-few-args")
		return
	}

	guild, err := p.extractGuild(event.Discord(), event.Fields()[2])
	if err != nil {
		event.Respond("whitelist.remove.invalid-invite")
		return
	}

	whitelistEntry, err := whitelistFind(p.db, "guild_id = ?", guild.ID)
	if err != nil {
		if !strings.Contains(err.Error(), "record not found") {
			event.Except(err)
			return
		}
	}

	if whitelistEntry == nil {
		event.Respond("whitelist.remove.not-found")
		return
	}

	if whitelistEntry.WhitelistedByUserID != event.UserID &&
		!event.Has(permissions.BotAdmin) {
		event.Respond("whitelist.remove.no-permissions")
		return
	}

	err = whitelistRemove(p.db, guild.ID)
	if err != nil {
		event.Except(err)
		return
	}

	p.logger.Info(
		"removed whitelisted server",
		zap.String("user_id", event.UserID),
		zap.String("guild_id", guild.ID),
	)

	err = cacheWhitelistAndBlacklist(p.db, p.redis)
	if err != nil {
		event.Except(err)
		return
	}

	_, err = event.Respond("whitelist.remove.success", "name", guild.Name)
	event.Except(err)
}
