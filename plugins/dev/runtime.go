package dev

import (
	"runtime"

	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleRuntime(event *events.Event) {
	_, err := event.Respond("dev.runtime.content",
		"goVersion", runtime.Version(),
		"goroutines", runtime.NumGoroutine(),
	)
	event.Except(err)
}
