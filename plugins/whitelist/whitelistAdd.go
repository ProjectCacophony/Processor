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

	blacklistEntry, err := blacklistFind(p.db, "guild_id = ?", guild.ID)
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

	whitelistEntry, err := whitelistFind(p.db, "guild_id = ?", guild.ID)
	if err != nil {
		if !strings.Contains(err.Error(), "record not found") {
			event.Except(err)
			return
		}
	}

	if whitelistEntry != nil {
		event.Respond("whitelist.add.already-whitelisted") // nolint: errcheck
		return
	}

	allUserEntries, err := whitelistFindMany(p.db,
		"whitelisted_by_user_id = ?", event.UserID,
	)
	if err != nil {
		event.Except(err)
		return
	}

	// TODO: bypassing for staff
	if len(allUserEntries) >= serversPerUserLimit {
		event.Respond("whitelist.add.too-many") // nolint: errcheck
		return
	}

	err = whitelistAdd(p.db, event.UserID, guild.ID)
	if err != nil {
		event.Except(err)
		return
	}

	p.logger.Info(
		"whitelisted server",
		zap.String("user_id", event.UserID),
		zap.String("guild_id", guild.ID),
	)

	err = cacheWhitelistAndBlacklist(p.db, p.redis)
	if err != nil {
		event.Except(err)
		return
	}

	_, err = event.Respond("whitelist.add.success", "name", guild.Name)
	event.Except(err)
}
