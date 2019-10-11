package dev

import (
	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleDM(event *events.Event) {
	if len(event.Fields()) < 3 {
		event.Except(events.NewUserError("common.to-few-params"))
		return
	}

	_, err := event.SendDM(event.UserID, event.Fields()[2])
	event.Except(err)

	event.React("ok")
}
