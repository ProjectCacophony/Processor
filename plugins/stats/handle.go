package stats

import (
	"gitlab.com/Cacophony/Processor/plugins/common"
	"gitlab.com/Cacophony/go-kit/events"
)

type Plugin struct {
}

func (p *Plugin) Name() string {
	return "stats"
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
		Description: "stats.help.description",
		Commands: []common.Command{
			{
				Name:        "stats.help.user.name",
				Description: "stats.help.user.description",
				Params: []common.CommandParam{
					{Name: "user", Type: common.Flag},
					{Name: "User", Type: common.User, Optional: true},
				},
			},
		},
	}
}

func (p *Plugin) Action(event *events.Event) bool {
	if !event.Command() {
		return false
	}

	if event.Fields()[0] != "stats" ||
		len(event.Fields()) < 2 {
		return false
	}

	if event.Fields()[1] == "user" {
		p.handleUser(event)
		return true
	}

	return false
}
