package dev

import (
	"strings"

	"gitlab.com/Cacophony/go-kit/config"
	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleConfig(event *events.Event) {
	if len(event.Fields()) < 4 {
		event.Respond("common.to-few-params")
		return
	}

	var err error
	var value []byte

	switch strings.ToLower(event.Fields()[2]) {
	case "guild":
		value, err = config.GuildGetBytes(event.DB(), event.GuildID, event.Fields()[3])
	case "user":
		value, err = config.UserGetBytes(event.DB(), event.GuildID, event.Fields()[3])
	default:
		event.Except(events.NewUserError("config type not found"))
		return
	}
	if err != nil {
		event.Except(err)
		return
	}

	if len(value) == 1 {
		switch value[0] {
		case 1:
			value = []byte("true")
		case 0:
			value = []byte("false")
		}
	}

	_, err = event.Respond("dev.config.content", "value", string(value))
	event.Except(err)
}
