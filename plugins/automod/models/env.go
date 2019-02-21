package models

import (
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/state"
)

type Env struct {
	Event   *events.Event
	State   *state.State
	GuildID string
	UserID  string
}
