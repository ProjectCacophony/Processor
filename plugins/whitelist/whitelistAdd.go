package whitelist

import (
	"strings"

	"go.uber.org/zap"

	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) whitelistAdd(event *events.Event) {
	if len(event.Fields()) < 2 {
		event.Respond("whitelist.add.too-few-args") // nolint: errcheck
		return
	}

	guild, err := p.extractGuild(event.Discord(), event.Fields()[1])
	if err != nil {
		event.Respond("whitelist.add.invalid-invite") // nolint: errcheck
		return
	}

	blacklistEntry, err := blacklistFindServer(p.db, guild.ID)
	if err != nil {
		if !strings.Contains(err.Error(), "record not found") {
			event.Except(err)
			return
		}
	}

	if blacklistEntry != nil {
		event.Respond("whitelist.add.blacklisted") // nolint: errcheck
		return
	}

	err = whitelistAddServer(p.db, event.UserID, guild.ID)
	if err != nil {
		if strings.Contains(err.Error(), "uix_whitelist_entries_guild_id") {
			event.Respond("whitelist.add.already-whitelisted") // nolint: errcheck
			return
		}
		event.Except(err)
	}

	p.logger.Info(
		"whitelisted server",
		zap.String("user_id", event.UserID),
		zap.String("guild_id", guild.ID),
	)

	_, err = event.Respond("whitelist.add.success", "name", guild.Name)
	event.Except(err)
}
