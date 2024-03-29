package instagram

import (
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/kelseyhightower/envconfig"
	"gitlab.com/Cacophony/Processor/plugins/common"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/interfaces"
	"gitlab.com/Cacophony/go-kit/permissions"
	"gitlab.com/Cacophony/go-kit/state"
	"go.uber.org/zap"
)

type Plugin struct {
	logger *zap.Logger
	state  *state.State
	db     *gorm.DB
}

func (p *Plugin) Names() []string {
	return []string{"instagram", "insta"}
}

type config struct {
	InstagramSessionIDs []string `envconfig:"INSTAGRAM_SESSION_IDS"`
}

func (p *Plugin) Start(params common.StartParameters) error {
	var config config
	err := envconfig.Process("", &config)
	if err != nil {
		return err
	}

	p.state = params.State
	p.logger = params.Logger
	p.db = params.DB

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
		Description: "instagram.help.description",
		Commands: []common.Command{
			{
				Name:        "instagram.help.list.name",
				Description: "instagram.help.list.description",
			},
			{
				Name:                "instagram.help.add.name",
				Description:         "instagram.help.add.description",
				PermissionsRequired: []interfaces.Permission{permissions.DiscordManageChannels},
				Params: []common.CommandParam{
					{Name: "add", Type: common.Flag},
					{Name: "Instagram Username", Type: common.Text},
					{Name: "channel", Type: common.Channel, Optional: true},
				},
			},
			{
				Name:                "instagram.help.remove.name",
				Description:         "instagram.help.remove.description",
				PermissionsRequired: []interfaces.Permission{permissions.DiscordManageChannels},
				Params: []common.CommandParam{
					{Name: "remove", Type: common.Flag},
					{Name: "Instagram Username", Type: common.Text},
				},
			},
			{
				Name:                "instagram.help.disable-post.name",
				Description:         "instagram.help.disable-post.description",
				PermissionsRequired: []interfaces.Permission{permissions.DiscordManageChannels},
				Params: []common.CommandParam{
					{Name: "disable-post", Type: common.Flag},
					{Name: "Instagram Username", Type: common.Text},
				},
			},
			{
				Name:                "instagram.help.enable-post.name",
				Description:         "instagram.help.enable-post.description",
				PermissionsRequired: []interfaces.Permission{permissions.DiscordManageChannels},
				Params: []common.CommandParam{
					{Name: "enable-post", Type: common.Flag},
					{Name: "Instagram Username", Type: common.Text},
				},
			},
			{
				Name:                "instagram.help.disable-story.name",
				Description:         "instagram.help.disable-story.description",
				PermissionsRequired: []interfaces.Permission{permissions.DiscordManageChannels, permissions.BotAdmin},
				Params: []common.CommandParam{
					{Name: "disable-story", Type: common.Flag},
					{Name: "Instagram Username", Type: common.Text},
				},
			},
			{
				Name:                "instagram.help.enable-story.name",
				Description:         "instagram.help.enable-story.description",
				PermissionsRequired: []interfaces.Permission{permissions.DiscordManageChannels, permissions.BotAdmin},
				Params: []common.CommandParam{
					{Name: "enable-story", Type: common.Flag},
					{Name: "Instagram Username", Type: common.Text},
				},
			},
		},
	}
}

func (p *Plugin) Action(event *events.Event) bool {
	if !event.Command() {
		return false
	}

	if event.Fields()[0] != "instagram" &&
		event.Fields()[0] != "insta" {
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

		case "disable-post", "disable-posts":
			event.RequireOr(func() {
				p.disable(event, modifyPosts)
			},
				permissions.DiscordManageChannels,
				permissions.DiscordChannelDM,
			)

			return true

		case "disable-story", "disable-stories":
			event.RequireOr(func() {
				p.disable(event, modifyStory)
			},
				permissions.DiscordManageChannels,
				permissions.DiscordChannelDM,
			)

			return true

		case "enable-post", "enable-posts":
			event.Require(func() {
				p.enable(event, modifyPosts)
			},
				permissions.BotAdmin,
				permissions.Or(
					permissions.DiscordManageChannels,
					permissions.DiscordChannelDM,
				),
			)

			return true

		case "enable-story", "enable-stories":
			event.Require(func() {
				p.enable(event, modifyStory)
			},
				permissions.BotAdmin,
				permissions.Or(
					permissions.DiscordManageChannels,
					permissions.DiscordChannelDM,
				),
			)

			return true

		}
	}

	p.status(event)

	return true
}
