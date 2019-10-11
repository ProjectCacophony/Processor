package dev

import (
	"errors"

	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleDevError(event *events.Event, user bool) {
	errorMessage := "undefined error"
	if len(event.Fields()) >= 3 {
		errorMessage = event.Fields()[2]
	}

	err := errors.New(errorMessage)
	if user {
		err = events.AsUserError(err)
	}

	event.Except(err)
}
