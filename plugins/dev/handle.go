package dev

import (
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/interfaces"
	"gitlab.com/Cacophony/go-kit/localisation"
)

type Plugin struct{}

func (p *Plugin) Name() string {
	return "dev"
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
	local, err := localisation.NewFileSource("assets/translations/dev.en.toml", "en")
	if err != nil {
		panic(err) // TODO: handle error
	}

	return []interfaces.Localisation{local}
}

func (p *Plugin) Action(event *events.Event) bool {
	if !event.Command() || event.Fields()[0] != "dev" {
		return false
	}

	if len(event.Fields()) < 2 {
		event.Respond("dev.no-subcommand") // nolint: errcheck
		return true
	}

	if event.Fields()[1] == "emoji" {

		handleDevEmoji(event)
		return true
	}

	event.Respond("dev.no-subcommand") // nolint: errcheck
	return true
}
