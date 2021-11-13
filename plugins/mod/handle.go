package mod

import (
	"gitlab.com/Cacophony/Processor/plugins/common"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/interfaces"
	"gitlab.com/Cacophony/go-kit/permissions"
)

type Plugin struct{}

func (p *Plugin) Names() []string {
	return []string{"mod"}
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
		Names:               p.Names(),
		Description:         "mod.help.description",
		PermissionsRequired: []interfaces.Permission{permissions.Not(permissions.DiscordChannelDM)},
		Commands: []common.Command{
			{
				Name:            "mod.help.moddm.name",
				Description:     "mod.help.moddm.description",
				SkipRootCommand: true,
				Params: []common.CommandParam{
					{Name: "mod-dm", Type: common.Flag},
					{Name: "User or User ID", Type: common.User},
					{Name: "Message Code", Type: common.Text},
				},
				PermissionsRequired: []interfaces.Permission{permissions.DiscordManageServer},
			},
			{
				Name:            "mod.help.modnote.name",
				Description:     "mod.help.modnote.description",
				SkipRootCommand: true,
				Params: []common.CommandParam{
					{Name: "mod-note", Type: common.Flag},
					{Name: "User or User ID", Type: common.User},
					{Name: "Message Code", Type: common.Text},
				},
				PermissionsRequired: []interfaces.Permission{permissions.DiscordManageServer},
			},
		},
	}
}

func (p *Plugin) Action(event *events.Event) bool {
	if !event.Command() {
		return false
	}
	if !event.Has(permissions.Not(permissions.DiscordChannelDM)) {
		return false
	}

	switch event.Fields()[0] {
	case "mod-dm", "moddm", "modm":
		event.Require(func() {
			p.handleModDM(event)
		}, permissions.DiscordManageServer)
		return true
	case "mod-note", "modnote":
		event.Require(func() {
			p.handleModNote(event)
		}, permissions.DiscordManageServer)
		return true
	}

	return false
}
