package stats

import (
	"gitlab.com/Cacophony/Processor/plugins/common"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/interfaces"
	"gitlab.com/Cacophony/go-kit/permissions"
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
			{
				Name:        "stats.help.server.name",
				Description: "stats.help.server.description",
				Params: []common.CommandParam{
					{Name: "server", Type: common.Flag},
				},
			},
			{
				Name:                "stats.help.server-specific.name",
				Description:         "stats.help.server-specific.description",
				PermissionsRequired: []interfaces.Permission{permissions.BotAdmin},
				Params: []common.CommandParam{
					{Name: "server", Type: common.Flag},
					{Name: "server ID", Type: common.Text},
				},
			},
			{
				Name:        "stats.help.channel.name",
				Description: "stats.help.channel.description",
				Params: []common.CommandParam{
					{Name: "channel", Type: common.Flag},
					{Name: "Channel", Type: common.Channel, Optional: true},
				},
			},
			{
				Name:        "stats.help.role.name",
				Description: "stats.help.role.description",
				Params: []common.CommandParam{
					{Name: "role", Type: common.Flag},
					{Name: "@role or role ID", Type: common.Text},
				},
			},
			{
				Name:                "stats.help.role-server.name",
				Description:         "stats.help.role-server.description",
				PermissionsRequired: []interfaces.Permission{permissions.BotAdmin},
				Params: []common.CommandParam{
					{Name: "role", Type: common.Flag},
					{Name: "@role or role ID", Type: common.Text},
					{Name: "Server ID", Type: common.Text},
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

	switch event.Fields()[1] {

	case "user", "member":
		p.handleUser(event)
		return true

	case "server", "guild":
		event.RequireOr(
			func() {
				p.handleServer(event)
			},
			permissions.Not(permissions.DiscordChannelDM),
			permissions.BotAdmin,
		)
		return true

	case "channel":
		p.handleChannel(event)
		return true

	case "role":
		event.RequireOr(
			func() {
				p.handleRole(event)
			},
			permissions.Not(permissions.DiscordChannelDM),
			permissions.BotAdmin,
		)
		return true

	}

	return false
}
