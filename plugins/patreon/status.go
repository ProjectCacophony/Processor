package patreon

import (
	"strings"

	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleStatus(event *events.Event) {
	// TOOD: use patreon permission to check qualification instead

	patron, err := getPatron(p.db, event.UserID)
	if err != nil && !strings.Contains(err.Error(), "record not found") {
		event.Except(err)
		return
	}

	var activePatron bool
	if patron != nil && patron.PatronStatus == "active_patron" {
		activePatron = true
	}

	_, err = event.Respond("patreon.status.response", "activePatron", activePatron)
	event.Except(err)
}
