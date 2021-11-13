package greeter

import (
	"gitlab.com/Cacophony/Processor/plugins/common"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/permissions"
)

type Plugin struct{}

func (p *Plugin) Names() []string {
	return []string{"greeter"}
}

func (p *Plugin) Start(params common.StartParameters) error {
	return params.DB.AutoMigrate(Entry{}).Error
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
		Description: "greeter.help.description",
		Commands: []common.Command{
			{
				Name:        "greeter.help.join.name",
				Description: "greeter.help.join.description",
				Params: []common.CommandParam{
					{Name: "join", Type: common.Flag},
					{Name: "channel", Type: common.Channel},
					{Name: "Message Code", Type: common.QuotedText, Optional: true},
					{Name: "Auto delete after (e.g. 5m)", Type: common.Duration, Optional: true},
				},
				PermissionsRequired: common.Permissions{permissions.DiscordManageChannels},
			},
			{
				Name:        "greeter.help.leave.name",
				Description: "greeter.help.leave.description",
				Params: []common.CommandParam{
					{Name: "leave", Type: common.Flag},
					{Name: "channel", Type: common.Channel},
					{Name: "Message Code", Type: common.QuotedText, Optional: true},
					{Name: "Auto delete after (e.g. 5m)", Type: common.Duration, Optional: true},
				},
				PermissionsRequired: common.Permissions{permissions.DiscordManageChannels},
			},
			{
				Name:        "greeter.help.ban.name",
				Description: "greeter.help.ban.description",
				Params: []common.CommandParam{
					{Name: "ban", Type: common.Flag},
					{Name: "channel", Type: common.Channel},
					{Name: "Message Code", Type: common.QuotedText, Optional: true},
					{Name: "Auto delete after (e.g. 5m)", Type: common.Duration, Optional: true},
				},
				PermissionsRequired: common.Permissions{permissions.DiscordManageChannels},
			},
			{
				Name:        "greeter.help.unban.name",
				Description: "greeter.help.unban.description",
				Params: []common.CommandParam{
					{Name: "unban", Type: common.Flag},
					{Name: "channel", Type: common.Channel},
					{Name: "Message Code", Type: common.QuotedText, Optional: true},
					{Name: "Auto delete after (e.g. 5m)", Type: common.Duration, Optional: true},
				},
				PermissionsRequired: common.Permissions{permissions.DiscordManageChannels},
			},
			{
				Name:        "greeter.help.status.name",
				Description: "greeter.help.status.description",
			},
		},
	}
}

func (p *Plugin) Action(event *events.Event) bool {
	if !event.Command() {
		return false
	}

	if event.Fields()[0] != "greeter" {
		return false
	}

	event.Require(func() {
		if len(event.Fields()) >= 2 {
			switch event.Fields()[1] {

			case "join":
				event.Require(func() {
					p.handleAdd(event, greeterTypeJoin)
				}, permissions.DiscordManageChannels)
				return

			case "leave":
				event.Require(func() {
					p.handleAdd(event, greeterTypeLeave)
				}, permissions.DiscordManageChannels)
				return

			case "ban":
				event.Require(func() {
					p.handleAdd(event, greeterTypeBan)
				}, permissions.DiscordManageChannels)
				return

			case "unban":
				event.Require(func() {
					p.handleAdd(event, greeterTypeUnban)
				}, permissions.DiscordManageChannels)
				return

			}
		}

		p.handleStatus(event)
	}, permissions.Not(permissions.DiscordChannelDM))

	return true
}
