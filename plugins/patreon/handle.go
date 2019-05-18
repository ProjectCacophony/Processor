package patreon

import (
	"gitlab.com/Cacophony/Processor/plugins/common"
	"gitlab.com/Cacophony/go-kit/events"
	"go.uber.org/zap"
)

type Plugin struct {
	logger *zap.Logger
}

func (p *Plugin) Name() string {
	return "patreon"
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

func (p *Plugin) Help() *common.PluginHelp {
	return &common.PluginHelp{
		Name:        p.Name(),
		Description: "patreon.help.description",
		Commands: []common.Command{{
			Name: "View Patreon Info",
		}},
	}
}

func (p *Plugin) Action(event *events.Event) bool {
	if !event.Command() {
		return false
	}

	if len(event.Fields()) < 1 || event.Fields()[0] != "help" {
		return false
	}

	// TODO: display supporters and patreon link in root command
	//   see robyul _patreon

	return false
}
