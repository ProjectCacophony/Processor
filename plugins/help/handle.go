package help

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
	return "help"
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
	local, err := localisation.NewFileSource("assets/translations/help.en.toml", "en")
	if err != nil {
		p.logger.Error("failed to load localisation", zap.Error(err))
	}

	return []interfaces.Localisation{local}
}

func (p *Plugin) Help() PluginHelp {
	return PluginHelp{
		PluginName: p.Name(),
		Hide:       true,
	}
}

func (p *Plugin) Action(event *events.Event) bool {
	if !event.Command() {
		return false
	}

	if len(event.Fields()) == 1 && event.Fields()[0] == "help" {
		listCommands(event, false)
	}

	if len(event.Fields()) > 1 && event.Fields()[1] == "public" {
		listCommands(event, true)
	}

	return false
}
