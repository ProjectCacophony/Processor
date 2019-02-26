package whitelist

import (
	"strings"

	"go.uber.org/zap"

	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) whitelistBlacklistAdd(event *events.Event) {
	if len(event.Fields()) < 3 {
		event.Respond("whitelist.blacklist-add.too-few-args") // nolint: errcheck
		return
	}

	guild, err := p.extractGuild(event.Discord(), event.Fields()[2])
	if err != nil {
		event.Respond("whitelist.blacklist-add.invalid-invite") // nolint: errcheck
		return
	}

	err = whitelistRemoveServer(p.db, guild.ID)
	if err != nil {
		event.Except(err)
		return
	}

	err = blacklistAddServer(p.db, event.UserID, guild.ID)
	if err != nil {
		if strings.Contains(err.Error(), "uix_whitelist_blacklist_entries_guild_id") {
			event.Respond("whitelist.blacklist-add.already-blacklisted") // nolint: errcheck
			return
		}
		event.Except(err)
		return
	}

	p.logger.Info(
		"blacklisted server",
		zap.String("user_id", event.UserID),
		zap.String("guild_id", guild.ID),
	)

	_, err = event.Respond("blacklist.add.success", "name", guild.Name)
	event.Except(err)
}
