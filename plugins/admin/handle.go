package admin

import (
	"gitlab.com/Cacophony/Processor/plugins/common"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/interfaces"
	"gitlab.com/Cacophony/go-kit/permissions"
	"gitlab.com/Cacophony/go-kit/state"
	"go.uber.org/zap"
)

type Plugin struct {
	state  *state.State
	logger *zap.Logger
}

func (p *Plugin) Name() string {
	return "sudo"
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

func (p *Plugin) Help() *common.PluginHelp {
	return &common.PluginHelp{
		Name:                p.Name(),
		Description:         "admin.help.description",
		PermissionsRequired: []interfaces.Permission{permissions.BotAdmin},
		Commands: []common.Command{
			{
				Name:        "admin.help.storage.toggle.user.name",
				Description: "admin.help.storage.toggle.user.description",
				Params: []common.CommandParam{
					{Name: "storage", Type: common.Flag},
					{Name: "enable/disable", Type: common.Flag},
					{Name: "user", Type: common.User},
				},
			},
		},
	}
}

func (p *Plugin) Action(event *events.Event) bool {
	if !event.Command() || event.Fields()[0] != "sudo" {
		return false
	}

	event.Require(func() {
		p.handleAsCommand(event)
	}, permissions.BotAdmin)
	return true
}

func (p *Plugin) handleAsCommand(event *events.Event) {
	if len(event.Fields()) < 2 {
		return
	}

	switch event.Fields()[1] {
	case "storage":
		if len(event.Fields()) < 3 {
			return
		}

		switch event.Fields()[2] {
		case "enable":
			p.toggleUserStorage(event, true)
			return
		case "disable":
			p.toggleUserStorage(event, false)

			return
		}

		return
	}
}
