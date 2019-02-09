package ping

import (
	"time"

	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/interfaces"
	"gitlab.com/Cacophony/go-kit/localisation"
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

func (p *Ping) Localisations() []interfaces.Localisation {
	local, err := localisation.NewFileSource("assets/translations/ping.en.toml", "en")
	if err != nil {
		panic(err) // TODO: handle error
	}

	return []interfaces.Localisation{local}
}

func (p *Ping) Action(event events.Event) bool {
	if event.Type != events.MessageCreateType {
		return false
	}

	switch event.MessageCreate.Content {
	case "ping":

		createdAt, err := event.MessageCreate.Timestamp.Parse()
		if err != nil {
			event.Except(err)
			return true
		}

		_, err = event.Sendf(
			event.MessageCreate.ChannelID,
			"ping.ping-response",
			"DiscordToGateway",
			event.ReceivedAt.Sub(createdAt),
			"GatewayToProcessor",
			time.Since(event.ReceivedAt),
		)
		if err != nil {
			event.Except(err)
			return true
		}

		return true

	case "pong":

		_, err := event.Send(event.MessageCreate.ChannelID, "ping.pong-response")
		if err != nil {
			event.Except(err)
			return true
		}

		return true
	}

	return false
}
