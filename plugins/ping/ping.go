package ping

import (
	"time"

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
	if event.Type != events.MessageCreateType {
		return false
	}

	if event.MessageCreate.Content != "ping" {
		return false
	}

	createdAt, err := event.MessageCreate.Timestamp.Parse()
	if err != nil {
		event.Except(err)
		return true
	}

	_, err = event.Discord().ChannelMessageSend(
		event.MessageCreate.ChannelID,
		"latency\ndiscord to gateway: "+event.ReceivedAt.Sub(createdAt).String()+"\n"+
			"gateway to processor: "+time.Since(event.ReceivedAt).String(),
	)
	if err != nil {
		event.Except(err)
		return true
	}

	return true
}
