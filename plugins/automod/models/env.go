package models

import (
	"encoding/json"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/state"
)

type HandlerInterface interface {
	Handle(event *events.Event) (process bool)
}

type Env struct {
	Rule      *Rule
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

func (e *Env) Marshal() ([]byte, error) {
	env := Env{
		Rule:      e.Rule,
		GuildID:   e.GuildID,
		ChannelID: e.ChannelID,
		UserID:    e.UserID,
		Messages:  e.Messages,
	}

	return json.Marshal(env)
}

func (e *Env) Unmarshal(data []byte) error {
	return json.Unmarshal(data, e)
}

type EnvMessage struct {
	ID        string
	ChannelID string
	Bot       bool
}

func NewEnvMessage(message *discordgo.Message) *EnvMessage {
	if message == nil {
		return nil
	}

	e := &EnvMessage{
		ID:        message.ID,
		ChannelID: message.ChannelID,
	}
	if message.Author != nil {
		e.Bot = message.Author.Bot
	}

	return e
}
