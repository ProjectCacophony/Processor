package tools

import (
	"gitlab.com/Cacophony/Processor/plugins/common"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/permissions"
)

type Plugin struct {
}

func (p *Plugin) Names() []string {
	return []string{"tools"}
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
		Description: "tools.help.description",
		Commands: []common.Command{
			{
				Name:            "tools.help.say.name",
				Description:     "tools.help.say.description",
				SkipRootCommand: true,
				Params: []common.CommandParam{
					{Name: "say", Type: common.Flag},
					{Name: "channel", Type: common.Channel},
					{Name: "message content, or code", Type: common.Text},
				},
				PermissionsRequired: common.Permissions{permissions.DiscordManageChannels},
			},
			{
				Name:            "tools.help.choose.name",
				Description:     "tools.help.choose.description",
				SkipRootCommand: true,
				Params: []common.CommandParam{
					{Name: "choose", Type: common.Flag},
					{Name: "item a", Type: common.QuotedText},
					{Name: "item b", Type: common.QuotedText},
					{Name: "…", Type: common.Text},
				},
			},
			{
				Name:            "tools.help.shorten.name",
				Description:     "tools.help.shorten.description",
				SkipRootCommand: true,
				Params: []common.CommandParam{
					{Name: "shorten", Type: common.Flag},
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

	switch event.Fields()[0] {
	case "shorten":
		p.handleShorten(event)
		return true
	case "choose":
		p.handleChoose(event)
		return true
	case "say":
		event.Require(func() {
			p.handleSay(event)
		}, permissions.DiscordManageChannels)
		return true
	}

	return false
}
