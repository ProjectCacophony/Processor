package shorten

import (
	"gitlab.com/Cacophony/Processor/plugins/common"
	"gitlab.com/Cacophony/go-kit/events"
)

type Plugin struct {
}

func (p *Plugin) Name() string {
	return "shorten"
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

func (p *Plugin) Help() *common.PluginHelp {
	return &common.PluginHelp{
		Name:        p.Name(),
		Description: "shorten.help.description",
		Commands: []common.Command{
			{
				Name:        "shorten.help.shorten.name",
				Description: "shorten.help.shorten.description",
				Params: []common.CommandParam{
					{Name: "link", Type: common.Link},
				},
			},
		},
	}
}

func (p *Plugin) Action(event *events.Event) bool {
	if !event.Command() {
		return false
	}

	p.handleShorten(event)

	return false
}
