package color

import (
	"gitlab.com/Cacophony/Processor/plugins/common"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/interfaces"
	"gitlab.com/Cacophony/go-kit/localisation"
	"go.uber.org/zap"
)

type Plugin struct {
	logger *zap.Logger
}

func (p *Plugin) Name() string {
	return "color"
}

func (p *Plugin) DBModels() []interface{} {
	return nil
}

func (p *Plugin) Start(params common.StartParameters) error {
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
	local, err := localisation.NewFileSource("assets/translations/color.en.toml", "en")
	if err != nil {
		p.logger.Error("failed to load localisation", zap.Error(err))
	}

	return []interfaces.Localisation{local}
}

func (p *Plugin) Action(event *events.Event) bool {
	if !event.Command() {
		return false
	}

	switch event.Fields()[0] {
	case "color", "colour":
		p.handleColor(event)

		return true

	}

	return false
}
