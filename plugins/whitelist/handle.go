package whitelist

import (
	"github.com/go-redis/redis"
	"gitlab.com/Cacophony/go-kit/interfaces"

	"gitlab.com/Cacophony/go-kit/permissions"

	"github.com/jinzhu/gorm"

	"gitlab.com/Cacophony/Processor/plugins/common"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/state"
	"go.uber.org/zap"
)

type Plugin struct {
	logger *zap.Logger
	state  *state.State
	db     *gorm.DB
	redis  *redis.Client
}

func (p *Plugin) Name() string {
	return "whitelist"
}

func (p *Plugin) Start(params common.StartParameters) error {
	var err error

	p.logger = params.Logger
	p.state = params.State
	p.db = params.DB
	p.redis = params.Redis

	err = params.DB.AutoMigrate(Entry{}, BlacklistEntry{}).Error
	if err != nil {
		return err
	}

	err = p.startWhitelistAndBlacklistCaching()
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
		Name:        p.Name(),
		Description: "whitelist.help.description",
		PermissionsRequired: []interfaces.Permission{
			permissions.BotAdmin,
			permissions.Patron,
		},
		Commands: []common.Command{
			{
				Name:        "whitelist.help.list.name",
				Description: "whitelist.help.list.description",
			},
			{
				Name:        "whitelist.help.add.name",
				Description: "whitelist.help.add.description",
				Params: []common.CommandParam{
					{Name: "add", Type: common.Flag},
					{Name: "Discord Invite", Type: common.DiscordInvite},
				},
			},
			{
				Name:        "whitelist.help.remove.name",
				Description: "whitelist.help.remove.description",
				Params: []common.CommandParam{
					{Name: "remove", Type: common.Flag},
					{Name: "Discord Invite", Type: common.DiscordInvite},
				},
			},
			{
				Name:                "whitelist.help.blacklist.name",
				Description:         "whitelist.help.blacklist.description",
				PermissionsRequired: []interfaces.Permission{permissions.BotAdmin},
				Params: []common.CommandParam{
					{Name: "blacklist", Type: common.Flag},
					{Name: "Discord Invite", Type: common.DiscordInvite},
				},
			},
		},
	}
}

func (p *Plugin) Action(event *events.Event) bool {
	if !event.Command() {
		return false
	}

	if event.Fields()[0] != "whitelist" {
		return false
	}

	if len(event.Fields()) >= 2 {

		switch event.Fields()[1] {
		case "list":

			p.whitelistList(event)
			return true

		case "remove":

			p.whitelistRemove(event)
			return true

		case "blacklist":

			event.Require(func() {
				p.blacklistAdd(event)
			},
				permissions.BotAdmin,
			)
			return true
		}

		event.RequireOr(func() {
			p.whitelistAdd(event)
		},
			permissions.BotAdmin,
			// TODO
			// permissions.Or(
			// 	permissions.BotAdmin,
			// 	permissions.Patron,
			// ),
		)
		return true
	}

	p.whitelistStatus(event)

	return true
}
