package instagram

import (
	"strings"

	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) remove(event *events.Event) {
	if len(event.Fields()) < 3 {
		// TODO: send message, too few args
		return
	}

	entry, err := entryFind(p.db,
		`
    ( ( (guild_id = ? AND dm = false) OR (channel_or_user_id = ? AND dm = true) ) AND dm = ? )
AND (instagram_username = ?)
`,
		event.GuildID, event.UserID, event.DM(), event.Fields()[2],
	)
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			event.Respond("instagram.remove.not-found")
			return
		}
		event.Except(err)
		return
	}

	err = entryRemove(p.db, entry.ID)
	if err != nil {
		event.Except(err)
	}

	_, err = event.Respond("instagram.remove.message", "entry", entry)
	event.Except(err)
}
