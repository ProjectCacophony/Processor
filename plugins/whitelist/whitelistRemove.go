package whitelist

import (
	"strings"

	"go.uber.org/zap"

	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) whitelistRemove(event *events.Event) {
	if len(event.Fields()) < 3 {
		event.Respond("whitelist.remove.too-few-args") // nolint: errcheck
		return
	}

	guild, err := p.extractGuild(event.Discord(), event.Fields()[2])
	if err != nil {
		event.Respond("whitelist.remove.invalid-invite") // nolint: errcheck
		return
	}

	whitelistEntry, err := whitelistFind(p.db, guild.ID)
	if err != nil {
		if !strings.Contains(err.Error(), "record not found") {
			event.Except(err)
			return
		}
	}

	if whitelistEntry == nil {
		event.Respond("whitelist.remove.not-found") // nolint: errcheck
		return
	}

	// TODO: bypassing for staff
	if whitelistEntry.WhitelistedByUserID != event.UserID {
		event.Respond("whitelist.remove.no-permissions") // nolint: errcheck
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

	_, err = event.Respond("whitelist.remove.success", "name", guild.Name)
	event.Except(err)
}
