package weverse

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/Seklfreak/geverse"
	"github.com/jinzhu/gorm"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"gitlab.com/Cacophony/Processor/plugins/common"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/interfaces"
	"gitlab.com/Cacophony/go-kit/permissions"
	"gitlab.com/Cacophony/go-kit/state"
)

type Plugin struct {
	state       *state.State
	db          *gorm.DB
	geverse     *geverse.Geverse
	communities []geverse.Community
}

func (p *Plugin) Names() []string {
	return []string{"weverse"}
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
		Names:       p.Names(),
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
			{
				Name:                "weverse.help.disable-artist.name",
				Description:         "weverse.help.disable-artist.description",
				PermissionsRequired: []interfaces.Permission{permissions.DiscordManageChannels},
				Params: []common.CommandParam{
					{Name: "disable-artist", Type: common.Flag},
					{Name: "Weverse community name", Type: common.Text},
				},
			},
			{
				Name:                "weverse.help.enable-artist.name",
				Description:         "weverse.help.enable-artist.description",
				PermissionsRequired: []interfaces.Permission{permissions.DiscordManageChannels},
				Params: []common.CommandParam{
					{Name: "enable-artist", Type: common.Flag},
					{Name: "Weverse community name", Type: common.Text},
				},
			},
			{
				Name:                "weverse.help.disable-media.name",
				Description:         "weverse.help.disable-media.description",
				PermissionsRequired: []interfaces.Permission{permissions.DiscordManageChannels},
				Params: []common.CommandParam{
					{Name: "disable-media", Type: common.Flag},
					{Name: "Weverse community name", Type: common.Text},
				},
			},
			{
				Name:                "weverse.help.enable-media.name",
				Description:         "weverse.help.enable-media.description",
				PermissionsRequired: []interfaces.Permission{permissions.DiscordManageChannels},
				Params: []common.CommandParam{
					{Name: "enable-media", Type: common.Flag},
					{Name: "Weverse community name", Type: common.Text},
				},
			},
			{
				Name:                "weverse.help.disable-notice.name",
				Description:         "weverse.help.disable-notice.description",
				PermissionsRequired: []interfaces.Permission{permissions.DiscordManageChannels},
				Params: []common.CommandParam{
					{Name: "disable-notice", Type: common.Flag},
					{Name: "Weverse community name", Type: common.Text},
				},
			},
			{
				Name:                "weverse.help.enable-notice.name",
				Description:         "weverse.help.enable-notice.description",
				PermissionsRequired: []interfaces.Permission{permissions.DiscordManageChannels},
				Params: []common.CommandParam{
					{Name: "enable-notice", Type: common.Flag},
					{Name: "Weverse community name", Type: common.Text},
				},
			},
			{
				Name:                "weverse.help.disable-moment.name",
				Description:         "weverse.help.disable-moment.description",
				PermissionsRequired: []interfaces.Permission{permissions.DiscordManageChannels},
				Params: []common.CommandParam{
					{Name: "disable-moment", Type: common.Flag},
					{Name: "Weverse community name", Type: common.Text},
				},
			},
			{
				Name:                "weverse.help.enable-moment.name",
				Description:         "weverse.help.enable-moment.description",
				PermissionsRequired: []interfaces.Permission{permissions.DiscordManageChannels},
				Params: []common.CommandParam{
					{Name: "enable-moment", Type: common.Flag},
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

		case "disable-artist", "disable-artists":
			event.RequireOr(func() {
				p.disable(event, modifyArtist)
			},
				permissions.DiscordManageChannels,
				permissions.DiscordChannelDM,
			)

			return true

		case "enable-artist", "enable-artists":
			event.RequireOr(func() {
				p.enable(event, modifyArtist)
			},
				permissions.DiscordManageChannels,
				permissions.DiscordChannelDM,
			)

			return true

		case "disable-media", "disable-medias":
			event.RequireOr(func() {
				p.disable(event, modifyMedia)
			},
				permissions.DiscordManageChannels,
				permissions.DiscordChannelDM,
			)

			return true

		case "enable-media", "enable-medias":
			event.RequireOr(func() {
				p.enable(event, modifyMedia)
			},
				permissions.DiscordManageChannels,
				permissions.DiscordChannelDM,
			)

			return true

		case "disable-notice", "disable-notices":
			event.RequireOr(func() {
				p.disable(event, modifyNotice)
			},
				permissions.DiscordManageChannels,
				permissions.DiscordChannelDM,
			)

			return true

		case "enable-notice", "enable-notices":
			event.RequireOr(func() {
				p.enable(event, modifyNotice)
			},
				permissions.DiscordManageChannels,
				permissions.DiscordChannelDM,
			)

			return true

		case "disable-moment", "disable-moments":
			event.RequireOr(func() {
				p.disable(event, modifyMoment)
			},
				permissions.DiscordManageChannels,
				permissions.DiscordChannelDM,
			)

			return true

		case "enable-moment", "enable-moments":
			event.RequireOr(func() {
				p.enable(event, modifyMoment)
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
