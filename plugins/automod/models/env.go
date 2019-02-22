package models

import (
	"github.com/go-redis/redis"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/state"
)

type Env struct {
	Event   *events.Event
	State   *state.State
	Redis   *redis.Client
	GuildID string
	UserID  string
}
