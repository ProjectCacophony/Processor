package color

import (
	"gitlab.com/Cacophony/Processor/plugins/common"
	"gitlab.com/Cacophony/go-kit/events"

	"go.uber.org/zap"
)

type Plugin struct {
	logger *zap.Logger
}

func (p *Plugin) Names() []string {
	return []string{"color", "colour"}
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
		Names:       p.Names(),
		Description: "color.help.description",
		Commands: []common.Command{{
			Params: []common.CommandParam{
				{Name: "Hex Code", Type: common.Text},
			},
		}, {
			Params: []common.CommandParam{
				{Name: "RGB Numbers", Type: common.Text},
			},
		}},
	}
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
