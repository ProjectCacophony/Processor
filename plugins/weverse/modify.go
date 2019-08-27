package weverse

import (
	"strings"

	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) disable(event *events.Event, modificationType modifyType) {
	if len(event.Fields()) < 3 {
		event.Respond("weverse.common.not-found")
		return
	}

	entry, err := entryFind(p.db,
		`
	    ( ( (guild_id = ? AND dm = false) OR (channel_or_user_id = ? AND dm = true) ) AND dm = ? )
	AND (LOWER(weverse_channel_name) = LOWER(?))
	`,
		event.GuildID, event.UserID, event.DM(), event.Fields()[2],
	)
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			event.Respond("weverse.common.not-found")
			return
		}
		event.Except(err)
		return
	}

	var alreadyApplied bool

	switch modificationType {
	case modifyArtist:
		if entry.DisableArtistFeed {
			alreadyApplied = true
		}
	case modifyMedia:
		if entry.DisableMediaFeed {
			alreadyApplied = true
		}
	case modifyNotice:
		if entry.DisableNoticeFeed {
			alreadyApplied = true
		}
	case modifyMoment:
		if entry.DisableMomentFeed {
			alreadyApplied = true
		}
	}

	if alreadyApplied {
		event.Respond("weverse.disable.already-applied")
		return
	}

	err = entryModify(p.db, entry.ID, modificationType, true)
	if err != nil {
		event.Except(err)
	}

	_, err = event.Respond("weverse.disable.success", "entry", entry, "type", modificationType)
	event.Except(err)
}

func (p *Plugin) enable(event *events.Event, modificationType modifyType) {
	if len(event.Fields()) < 3 {
		event.Respond("weverse.common.not-found")
		return
	}

	entry, err := entryFind(p.db,
		`
	    ( ( (guild_id = ? AND dm = false) OR (channel_or_user_id = ? AND dm = true) ) AND dm = ? )
	AND (LOWER(weverse_channel_name) = LOWER(?))
	`,
		event.GuildID, event.UserID, event.DM(), event.Fields()[2],
	)
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			event.Respond("weverse.common.not-found")
			return
		}
		event.Except(err)
		return
	}

	var alreadyApplied bool

	switch modificationType {
	case modifyArtist:
		if !entry.DisableArtistFeed {
			alreadyApplied = true
		}
	case modifyMedia:
		if !entry.DisableMediaFeed {
			alreadyApplied = true
		}
	case modifyNotice:
		if !entry.DisableNoticeFeed {
			alreadyApplied = true
		}
	case modifyMoment:
		if !entry.DisableMomentFeed {
			alreadyApplied = true
		}
	}

	if alreadyApplied {
		event.Respond("weverse.enable.already-applied")
		return
	}

	err = entryModify(p.db, entry.ID, modificationType, false)
	if err != nil {
		event.Except(err)
	}

	_, err = event.Respond("weverse.enable.success", "entry", entry, "type", modificationType)
	event.Except(err)
}
