package vlive

import (
	"strings"

	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) remove(event *events.Event) {
	if len(event.Fields()) < 3 {
		event.Respond("vlive.common.not-found")
		return
	}

	entry, err := entryFind(p.db,
		`
    ( ( (guild_id = ? AND dm = false) OR (channel_or_user_id = ? AND dm = true) ) AND dm = ? )
AND (v_live_channel_id = ?)
`,
		event.GuildID, event.UserID, event.DM(), extractVLiveChannelID(event.Fields()[2]),
	)
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			event.Respond("vlive.remove.not-found")
			return
		}
		event.Except(err)
		return
	}

	err = entryRemove(p.db, entry.ID)
	if err != nil {
		event.Except(err)
	}

	_, err = event.Respond("vlive.remove.message", "entry", entry)
	event.Except(err)
}
