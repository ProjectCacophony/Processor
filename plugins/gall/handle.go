package gall

import (
	"net/http"
	"strings"
	"time"

	"gitlab.com/Cacophony/go-kit/interfaces"
	"gitlab.com/Cacophony/go-kit/permissions"

	"github.com/Seklfreak/ginside"
	"github.com/jinzhu/gorm"
	"gitlab.com/Cacophony/go-kit/state"

	"gitlab.com/Cacophony/Processor/plugins/common"
	"gitlab.com/Cacophony/go-kit/events"
	"go.uber.org/zap"
)

type Plugin struct {
	logger  *zap.Logger
	state   *state.State
	db      *gorm.DB
	ginside *ginside.GInside
}

func (p *Plugin) Names() []string {
	return []string{"gall"}
}

func (p *Plugin) Start(params common.StartParameters) error {
	var err error

	p.state = params.State
	p.logger = params.Logger
	p.db = params.DB
	p.ginside = ginside.NewGInside(&http.Client{
		Timeout: time.Minute,
	})

	err = params.DB.AutoMigrate(Entry{}).Error
	if err != nil {
		return err
	}

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
		Description: "gall.help.description",
		Commands: []common.Command{
			{
				Name:        "gall.help.list.name",
				Description: "gall.help.list.description",
			},
			{
				Name:                "gall.help.add.name",
				Description:         "gall.help.add.description",
				PermissionsRequired: []interfaces.Permission{permissions.DiscordManageChannels},
				Params: []common.CommandParam{
					{Name: "add", Type: common.Flag},
					{Name: "Gallery Names", Type: common.Text},
					{Name: "channel", Type: common.Channel, Optional: true},
					{Name: "all", Type: common.Flag, Optional: true},
				},
			},
			{
				Name:                "gall.help.remove.name",
				Description:         "gall.help.remove.description",
				PermissionsRequired: []interfaces.Permission{permissions.DiscordManageChannels},
				Params: []common.CommandParam{
					{Name: "remove", Type: common.Flag},
					{Name: "Gallery Names", Type: common.Text},
				},
			},
		},
	}
}

func (p *Plugin) Action(event *events.Event) bool {
	if !event.Command() {
		return false
	}

	if event.Fields()[0] != "gall" {
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
