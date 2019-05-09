package whitelist

import (
	"regexp"

	"github.com/go-redis/redis"

	"gitlab.com/Cacophony/go-kit/permissions"

	"github.com/jinzhu/gorm"

	"gitlab.com/Cacophony/Processor/plugins/common"
	"gitlab.com/Cacophony/Processor/plugins/help"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/interfaces"
	"gitlab.com/Cacophony/go-kit/localisation"
	"gitlab.com/Cacophony/go-kit/state"
	"go.uber.org/zap"
)

type Plugin struct {
	logger *zap.Logger
	state  *state.State
	db     *gorm.DB
	redis  *redis.Client

	discordInviteRegex *regexp.Regexp
	snowflakeRegex     *regexp.Regexp
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

	p.discordInviteRegex, err = regexp.Compile(
		`^(http(s)?:\/\/)?(discord\.gg(\/invite)?|discordapp\.com\/invite)\/([A-Za-z0-9-]+)$`,
	)
	if err != nil {
		return err
	}

	p.snowflakeRegex, err = regexp.Compile(
		`^[0-9]+$`,
	)
	if err != nil {
		return err
	}

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

func (p *Plugin) Localisations() []interfaces.Localisation {
	local, err := localisation.NewFileSource("assets/translations/whitelist.en.toml", "en")
	if err != nil {
		p.logger.Error("failed to load localisation", zap.Error(err))
	}

	return []interfaces.Localisation{local}
}

func (p *Plugin) Help() help.PluginHelp {
	return help.PluginHelp{
		Name: p.Name(),
		Hide: true,
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

			event.Require(func() {
				p.whitelistList(event)
			}, permissions.BotAdmin)
			return true

		case "remove":

			p.whitelistRemove(event)
			return true

		case "blacklist":

			event.Require(func() {
				p.blacklistAdd(event)
			}, permissions.BotAdmin)
			return true
		}

		event.RequireOr(func() {
			p.whitelistAdd(event)
		},
			permissions.BotAdmin,
		)
		return true
	}

	p.whitelistStatus(event)

	return true
}
