package ping

import (
	"time"

	"gitlab.com/Cacophony/Processor/pkg/kit"
	"gitlab.com/Cacophony/go-kit/events"
)

type Ping struct{}

func (p *Ping) Name() string {
	return "ping"
}

func (p *Ping) Start() error {
	return nil
}

func (p *Ping) Stop() error {
	return nil
}

func (p *Ping) Priority() int {
	return 0
}

func (p *Ping) Passthrough() bool {
	return false
}

func (p *Ping) Action(event events.Event) bool {
	if event.Type != events.MessageCreateEventType {
		return false
	}

	if event.MessageCreate.Content != "ping" {
		return false
	}

	createdAt, _ := event.MessageCreate.Timestamp.Parse()

	session, _ := kit.BotSession(event.BotUserID)

	_, _ = session.ChannelMessageSend(
		event.MessageCreate.ChannelID,
		"latency\ndiscord to gateway: "+event.ReceivedAt.Sub(createdAt).String()+"\n"+
			"gateway to processor: "+time.Since(event.ReceivedAt).String(),
	)

	return true
}
