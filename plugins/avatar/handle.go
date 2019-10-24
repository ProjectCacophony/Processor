package avatar

import (
	"gitlab.com/Cacophony/Processor/plugins/common"
	"gitlab.com/Cacophony/go-kit/events"
)

type Plugin struct{}

func (p *Plugin) Names() []string {
	return []string{"avatar", "av"}
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
		Names:       p.Names(),
		Description: "avatar.help.description",
		Commands: []common.Command{
			{
				Name:        "avatar.help.avatar",
				Description: "avatar.help.avatar.description",
				Params: []common.CommandParam{
					{Name: "User", Type: common.User, Optional: true},
				},
			},
			{
				Name:        "avatar.help.size-avatar",
				Description: "avatar.help.size-avatar.description",
				Params: []common.CommandParam{
					{Name: "size", Type: common.Flag},
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

	switch event.Fields()[0] {
	case "avatar", "av":
		p.displayUserAvatar(event)
		return true
	}

	return false
}
