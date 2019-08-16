package weverse

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/Seklfreak/geverse"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"gitlab.com/Cacophony/go-kit/interfaces"
	"gitlab.com/Cacophony/go-kit/permissions"

	"github.com/jinzhu/gorm"
	"gitlab.com/Cacophony/go-kit/state"

	"gitlab.com/Cacophony/Processor/plugins/common"
	"gitlab.com/Cacophony/go-kit/events"
	"go.uber.org/zap"
)

type Plugin struct {
	logger      *zap.Logger
	state       *state.State
	db          *gorm.DB
	geverse     *geverse.Geverse
	communities []geverse.Community
}

func (p *Plugin) Name() string {
	return "weverse"
}

type weverseConfig struct {
	WeverseToken string `envconfig:"WEVERSE_TOKEN"`
}

func (p *Plugin) Start(params common.StartParameters) error {
	var config weverseConfig
	err := envconfig.Process("", &config)
	if err != nil {
		return errors.Wrap(err, "failure loading weverse module config")
	}

	p.state = params.State
	p.logger = params.Logger
	p.db = params.DB
	p.geverse = geverse.NewGeverse(
		&http.Client{
			Timeout: 30 * time.Second,
		},
		config.WeverseToken,
	)

	meInfo, err := p.geverse.GetMe(context.Background())
	if err != nil {
		return err
	}
	p.communities = meInfo.MyCommunities

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
		Description: "weverse.help.description",
		Commands: []common.Command{
			{
				Name:        "weverse.help.list.name",
				Description: "weverse.help.list.description",
			},
			{
				Name:                "weverse.help.add.name",
				Description:         "weverse.help.add.description",
				PermissionsRequired: []interfaces.Permission{permissions.DiscordManageChannels},
				Params: []common.CommandParam{
					{Name: "add", Type: common.Flag},
					{Name: "Weverse community name", Type: common.Text},
					{Name: "channel", Type: common.Channel, Optional: true},
				},
			},
			{
				Name:                "weverse.help.remove.name",
				Description:         "weverse.help.remove.description",
				PermissionsRequired: []interfaces.Permission{permissions.DiscordManageChannels},
				Params: []common.CommandParam{
					{Name: "remove", Type: common.Flag},
					{Name: "Weverse community name", Type: common.Text},
				},
			},
		},
	}
}

func (p *Plugin) Action(event *events.Event) bool {
	if !event.Command() {
		return false
	}

	if event.Fields()[0] != "weverse" {
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
