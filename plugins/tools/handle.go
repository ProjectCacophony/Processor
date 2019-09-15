package tools

import (
	"math/rand"
	"time"

	"gitlab.com/Cacophony/Processor/plugins/common"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/permissions"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

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
				PermissionsRequired: common.Permissions{
					permissions.DiscordManageChannels,
					permissions.DiscordManageMessages,
				},
			},
			{
				Name:            "tools.help.edit.name",
				Description:     "tools.help.edit.description",
				SkipRootCommand: true,
				Params: []common.CommandParam{
					{Name: "edit", Type: common.Flag},
					{Name: "message link", Type: common.Link},
					{Name: "message content, or code", Type: common.Text},
				},
				PermissionsRequired: common.Permissions{
					permissions.DiscordManageChannels,
					permissions.DiscordManageMessages,
				},
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
				Name:            "tools.help.roll.name",
				Description:     "tools.help.roll.description",
				SkipRootCommand: true,
				Params: []common.CommandParam{
					{Name: "roll", Type: common.Flag},
					{Name: "maximum number", Type: common.Text},
				},
			},
			{
				Name:            "tools.help.8ball.name",
				Description:     "tools.help.8ball.description",
				SkipRootCommand: true,
				Params: []common.CommandParam{
					{Name: "8ball", Type: common.Flag},
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
	case "roll":
		p.handleRoll(event)
		return true
	case "8ball":
		p.handle8ball(event)
		return true
	case "say":
		event.Require(func() {
			p.handleSay(event)
		}, permissions.Or(permissions.DiscordManageChannels, permissions.DiscordManageMessages))
		return true
	case "edit":
		event.Require(func() {
			p.handleEdit(event)
		}, permissions.DiscordManageChannels, permissions.DiscordManageMessages)
		return true
	}

	return false
}
