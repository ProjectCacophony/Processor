package help

import (
	"gitlab.com/Cacophony/Processor/plugins/common"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/interfaces"
	"gitlab.com/Cacophony/go-kit/localization"
	"go.uber.org/zap"
)

type Plugin struct {
	logger         *zap.Logger
	pluginHelpList []*common.PluginHelp
}

func (p *Plugin) Name() string {
	return "help"
}

func (p *Plugin) Start(params common.StartParameters) error {
	p.logger = params.Logger
	p.pluginHelpList = params.PluginHelpList

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

func (p *Plugin) Localizations() []interfaces.Localization {
	local, err := localization.NewFileSource("assets/translations/help.en.toml", "en")
	if err != nil {
		p.logger.Error("failed to load localization", zap.Error(err))
	}

	return []interfaces.Localization{local}
}

func (p *Plugin) Help() *common.PluginHelp {
	return &common.PluginHelp{
		Name: p.Name(),
		Hide: true,
	}
}

func (p *Plugin) Action(event *events.Event) bool {
	if !event.Command() {
		return false
	}

	if len(event.Fields()) == 1 && event.Fields()[0] == "help" {
		listCommands(event, p.pluginHelpList, false)
	}

	if len(event.Fields()) > 1 && event.Fields()[1] == "public" {
		listCommands(event, p.pluginHelpList, true)
	}

	return false
}
