package dev

import (
	"gitlab.com/Cacophony/Processor/plugins/common"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/interfaces"
	"gitlab.com/Cacophony/go-kit/localisation"
	"gitlab.com/Cacophony/go-kit/state"
	"go.uber.org/zap"
)

type Plugin struct {
	state  *state.State
	logger *zap.Logger
}

func (p *Plugin) Name() string {
	return "dev"
}

func (p *Plugin) DBModels() []interface{} {
	return nil
}

func (p *Plugin) Start(params common.StartParameters) error {
	p.state = params.State
	p.logger = params.Logger
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
		p.logger.Error("failed to load localisation", zap.Error(err))
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

		p.handleDevEmoji(event)
		return true
	case "sleep":

		p.handleDevSleep(event)
		return true
	case "state":
		if len(event.Fields()) > 2 {
			if event.Fields()[2] == "guilds" {
				p.handleDevStateGuilds(event)
				return true
			}
		}

		p.handleDevState(event)
		return true
	case "translate":

		p.handleDevTranslate(event)
		return true
	case "error":

		p.handleDevError(event)
		return true
	}

	event.Respond("dev.no-subcommand") // nolint: errcheck
	return true
}
