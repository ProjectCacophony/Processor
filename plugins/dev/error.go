package dev

import (
	"errors"

	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleDevError(event *events.Event) {
	errorMessage := "undefined error"
	if len(event.Fields()) >= 3 {
		errorMessage = event.Fields()[2]
	}

	event.Except(errors.New(errorMessage))
}
