package admin

import (
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/permissions"
)

func (p *Plugin) toggleUserStorage(event *events.Event, enable bool) {
	if len(event.Fields()) < 4 {
		event.Respond("common.to-few-params")
		return
	}

	// check if user exists
	targetUser, err := event.State().UserFromMention(event.Fields()[3])
	if err != nil {
		event.Except(err)
		return
	}

	// give/remove user permission
	if enable {
		err = permissions.CacoFileStorage.Give(event.DB(), targetUser.ID)
	} else {
		err = permissions.CacoFileStorage.Remove(event.DB(), targetUser.ID)
	}
	if err != nil {
		event.Except(err)
		return
	}

	event.Respond("admin.storage.toggle",
		"enable", enable,
		"permission", permissions.CacoFileStorage.Name(),
		"username", targetUser.Username,
	)
}
