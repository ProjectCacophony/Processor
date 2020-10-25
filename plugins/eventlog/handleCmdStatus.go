package eventlog

import (
	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleCmdStatus(event *events.Event) {
	_, err := event.Respond("eventlog.status.message", "enabled", isEnabled(event))
	event.Except(err)
}
