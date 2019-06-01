package dev

import (
	"strconv"

	"gitlab.com/Cacophony/go-kit/permissions"

	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleDevPermission(event *events.Event) {
	if len(event.Fields()) < 3 {
		return
	}

	userID := event.UserID
	if len(event.MessageCreate.Mentions) > 0 {
		userID = event.MessageCreate.Mentions[0].ID
	}

	permissionID, err := strconv.Atoi(event.Fields()[2])
	if err != nil {
		event.Except(err)
		return
	}

	has := permissions.NewDiscordPermission("", permissionID).Match(
		p.state,
		event.DB(),
		userID,
		event.ChannelID,
		event.DM(),
	)

	_, err = event.Respond("dev.permission",
		"has", has, "permissionID", permissionID, "userID", userID)
	if err != nil {
		event.Except(err)
		return
	}
}
