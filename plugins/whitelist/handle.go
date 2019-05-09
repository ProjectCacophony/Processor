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

	// nolint: gocritic
	p.discordInviteRegex, err = regexp.Compile(
		`^(http(s)?:\/\/)?(discord\.gg(\/invite)?|discordapp\.com\/invite)\/([A-Za-z0-9-]+)$`,
	)
	if err != nil {
		return err
	}

	// nolint: gocritic
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
			}, permissions.BotOwner)
			return true

		case "remove":

			p.whitelistRemove(event)
			return true

		case "blacklist":

			event.Require(func() {
				p.blacklistAdd(event)
			}, permissions.BotOwner)
			return true
		}

		event.RequireOr(func() {
			p.whitelistAdd(event)
		},
			permissions.BotOwner,
			// TODO: don't hardcode
			// sekl's dev cord / Cacophony Whitelist Role
			permissions.NewDiscordRole(
				p.state, "208673735580844032", "549951245738180688",
			),
			// Cacophony / Cacophony Team
			permissions.NewDiscordRole(
				p.state, "435420687906111498", "440519691904090113",
			),
		)
		return true
	}

	p.whitelistStatus(event)

	return true
}
