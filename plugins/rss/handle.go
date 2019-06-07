package rss

import (
	"strings"

	"github.com/mmcdole/gofeed"

	"gitlab.com/Cacophony/go-kit/permissions"

	"github.com/jinzhu/gorm"
	"gitlab.com/Cacophony/go-kit/state"

	"gitlab.com/Cacophony/Processor/plugins/common"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/interfaces"
	"go.uber.org/zap"
)

type Plugin struct {
	logger *zap.Logger
	state  *state.State
	db     *gorm.DB
	parser *gofeed.Parser
}

func (p *Plugin) Name() string {
	return "rss"
}

func (p *Plugin) Start(params common.StartParameters) error {
	p.state = params.State
	p.logger = params.Logger
	p.db = params.DB
	p.parser = gofeed.NewParser()

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
		Name:        p.Name(),
		Description: "rss.help.description",
		Commands: []common.Command{
			{
				Name:        "rss.help.list.name",
				Description: "rss.help.list.description",
			},
			{
				Name:                "rss.help.add.name",
				Description:         "rss.help.add.description",
				PermissionsRequired: []interfaces.Permission{permissions.DiscordManageChannels},
				Params: []common.CommandParam{
					{Name: "add", Type: common.Flag},
					{Name: "link", Type: common.Text},
					{Name: "channel", Type: common.Channel, Optional: true},
				},
			}, {
				Name:                "rss.help.add.description",
				Description:         "rss.help.remove.description",
				PermissionsRequired: []interfaces.Permission{permissions.DiscordManageChannels},
				Params: []common.CommandParam{
					{Name: "remove", Type: common.Flag},
					{Name: "link", Type: common.Text},
				},
			},
		},
	}
}

func (p *Plugin) Action(event *events.Event) bool {
	if !event.Command() {
		return false
	}

	if event.Fields()[0] != "rss" {
		return false
	}

	if len(event.Fields()) > 1 {
		switch strings.ToLower(event.Fields()[1]) {
		case "add":
			event.RequireOr(func() {
				p.add(event)
			},
				permissions.DiscordManageChannels,
				permissions.DiscordChannelDM,
			)

			return true

		case "remove":
			event.RequireOr(func() {
				p.remove(event)
			},
				permissions.DiscordManageChannels,
				permissions.DiscordChannelDM,
			)

			return true
		}
	}

	p.status(event)

	return true
}
