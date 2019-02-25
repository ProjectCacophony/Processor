package models

import (
	"github.com/go-redis/redis"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/state"
)

type HandlerInterface interface {
	Handle(event *events.Event) (process bool)
}

type Env struct {
	Event     *events.Event
	State     *state.State
	Redis     *redis.Client
	Handler   HandlerInterface
	GuildID   string
	ChannelID []string
	UserID    []string
	Messages  []*EnvMessage
	Tokens    map[string]string
}

type EnvMessage struct {
	ID       string
	ChanneID string
}
