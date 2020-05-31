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
	state     *state.State
	logger    *zap.Logger
	publisher *events.Publisher
}

func (p *Plugin) Names() []string {
	return []string{"sudo"}
}

func (p *Plugin) Start(params common.StartParameters) error {
	p.state = params.State
	p.logger = params.Logger
	p.publisher = params.Publisher
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
		Names:               p.Names(),
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
			{
				Name:        "admin.help.as.name",
				Description: "admin.help.as.description",
				Params: []common.CommandParam{
					{Name: "as", Type: common.Flag},
					{Name: "user", Type: common.User},
					{Name: "command", Type: common.Text},
					{Name: "…", Type: common.Text},
				},
			},
			{
				Name:        "admin.help.do.name",
				Description: "admin.help.do.description",
				Params: []common.CommandParam{
					{Name: "do", Type: common.Flag},
					{Name: "command", Type: common.Text},
					{Name: "…", Type: common.Text},
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

	case "as":
		if len(event.Fields()) < 3 {
			return
		}

		event.Require(func() {
			p.handleAs(event)
		}, permissions.Not(permissions.DiscordChannelDM))
		return

	case "do":
		if len(event.Fields()) < 3 {
			return
		}

		event.Require(func() {
			p.handleDo(event)
		}, permissions.Not(permissions.DiscordChannelDM))
		return
	}
}
