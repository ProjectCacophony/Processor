package dev

import (
	"gitlab.com/Cacophony/Processor/plugins/common"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/interfaces"
	"gitlab.com/Cacophony/go-kit/localisation"
)

type Plugin struct{}

func (p *Plugin) Name() string {
	return "dev"
}

func (p *Plugin) DBModels() []interface{} {
	return nil
}

func (p *Plugin) Start(params common.StartParameters) error {
	return nil
}

func (p *Plugin) Stop(params common.StopParameters) error {
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

	switch event.Fields()[1] {
	case "emoji":

		handleDevEmoji(event)
		return true
	case "sleep":

		handleDevSleep(event)
		return true
	case "state":

		handleDevState(event)
		return true
	case "translate":

		handleDevTranslate(event)
		return true
	}

	event.Respond("dev.no-subcommand") // nolint: errcheck
	return true
}
