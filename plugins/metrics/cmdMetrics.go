package metrics

import (
	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleCmdMetrics(event *events.Event) {
	_, err := event.Respond(
		"metrics.content",
		"metrics", metrics,
		"db", event.DB(),
	)
	event.Except(err)
}
