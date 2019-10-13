package eventlog

import (
	"gitlab.com/Cacophony/go-kit/events"
)

type enhancedItem struct {
	Item
	Summary string
}

func (p *Plugin) handleCmdHistory(event *events.Event) {
	targetUser, err := event.FindUser(events.WithoutFallbackToSelf())
	if err != nil {
		event.Except(err)
		return
	}

	items, err := FindManyItem(event.DB(), 10,
		"author_id = ? OR (target_type = ? AND target_value = ?)", targetUser.ID, EntityTypeUser, targetUser.ID,
	)
	if err != nil {
		event.Except(err)
		return
	}

	enhancedItems := make([]enhancedItem, len(items))
	for i, item := range items {
		enhancedItems[len(items)-1-i].Item = item
		enhancedItems[len(items)-1-i].Summary = item.Summary(event.State(), targetUser.ID)
	}

	_, err = event.Respond("eventlog.history.content", "items", enhancedItems, "user", targetUser)
	event.Except(err)
}
