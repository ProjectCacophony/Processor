package ping

import (
	"time"

	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/interfaces"
	"gitlab.com/Cacophony/go-kit/localisation"
)

type Plugin struct{}

func (p *Plugin) Name() string {
	return "ping"
}

func (p *Plugin) Start() error {
	return nil
}

func (p *Plugin) Stop() error {
	return nil
}

func (p *Plugin) Priority() int {
	return 0
}

func (p *Plugin) Passthrough() bool {
	return false
}

func (p *Plugin) Localisations() []interfaces.Localisation {
	local, err := localisation.NewFileSource("assets/translations/ping.en.toml", "en")
	if err != nil {
		panic(err) // TODO: handle error
	}

	return []interfaces.Localisation{local}
}

func (p *Plugin) Action(event *events.Event) bool {
	if !event.Command() {
		return false
	}

	switch event.Fields()[0] {
	case "ping":

		createdAt, err := event.MessageCreate.Timestamp.Parse()
		if err != nil {
			event.Except(err)
			return true
		}

		_, err = event.Respondf(
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

		_, err := event.Respond("ping.pong-response")
		if err != nil {
			event.Except(err)
			return true
		}

		return true
	}

	return false
}
