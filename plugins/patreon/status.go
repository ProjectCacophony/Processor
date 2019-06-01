package patreon

import (
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/permissions"
)

func (p *Plugin) handleStatus(event *events.Event) {
	var activePatron bool
	if event.Has(permissions.Patron) {
		activePatron = true
	}

	_, err := event.Respond(
		"patreon.status.response",
		"activePatron", activePatron,
	)
	event.Except(err)
}
