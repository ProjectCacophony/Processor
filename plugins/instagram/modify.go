package instagram

import (
	"strings"

	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) disable(event *events.Event, disableType modifyType) {
	if len(event.Fields()) < 3 {
		event.Respond("instagram.common.not-found")
		return
	}

	entry, err := entryFind(p.db,
		`
    ( ( (guild_id = ? AND dm = false) OR (channel_or_user_id = ? AND dm = true) ) AND dm = ? )
AND (instagram_username = ?)
`,
		event.GuildID, event.UserID, event.DM(), extractInstagramUsername(event.Fields()[2]),
	)
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			event.Respond("instagram.common.not-found")
			return
		}
		event.Except(err)
		return
	}

	var alreadyApplied bool

	switch disableType {
	case modifyPosts:
		if entry.DisablePostFeed {
			alreadyApplied = true
		}
	case modifyStory:
		if entry.DisableStoryFeed {
			alreadyApplied = true
		}
	}

	if alreadyApplied {
		event.Respond("instagram.disable.already-applied")
		return
	}

	err = entryModify(p.db, entry.ID, disableType, true)
	if err != nil {
		event.Except(err)
	}

	_, err = event.Respond("instagram.disable.success", "entry", entry, "type", disableType)
	event.Except(err)
}

func (p *Plugin) enable(event *events.Event, enableType modifyType) {
	if len(event.Fields()) < 3 {
		event.Respond("instagram.common.not-found")
		return
	}

	entry, err := entryFind(p.db,
		`
    ( ( (guild_id = ? AND dm = false) OR (channel_or_user_id = ? AND dm = true) ) AND dm = ? )
AND (instagram_username = ?)
`,
		event.GuildID, event.UserID, event.DM(), extractInstagramUsername(event.Fields()[2]),
	)
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			event.Respond("instagram.common.not-found")
			return
		}
		event.Except(err)
		return
	}

	var alreadyApplied bool

	switch enableType {
	case modifyPosts:
		if !entry.DisablePostFeed {
			alreadyApplied = true
		}
	case modifyStory:
		if !entry.DisableStoryFeed {
			alreadyApplied = true
		}
	}

	if alreadyApplied {
		event.Respond("instagram.enable.already-applied")
		return
	}

	err = entryModify(p.db, entry.ID, enableType, false)
	if err != nil {
		event.Except(err)
	}

	_, err = event.Respond("instagram.enable.success", "entry", entry, "type", enableType)
	event.Except(err)
}
