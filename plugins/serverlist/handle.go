package serverlist

import (
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"gitlab.com/Cacophony/Processor/plugins/common"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/interfaces"
	"gitlab.com/Cacophony/go-kit/localisation"
	"gitlab.com/Cacophony/go-kit/permissions"
	"gitlab.com/Cacophony/go-kit/state"
	"go.uber.org/zap"
)

type Plugin struct {
	logger *zap.Logger
	db     *gorm.DB
	state  *state.State
	redis  *redis.Client
	tokens map[string]string
}

func (p *Plugin) Name() string {
	return "serverlist"
}

func (p *Plugin) Start(params common.StartParameters) error {
	var err error

	p.logger = params.Logger
	p.db = params.DB
	p.state = params.State
	p.redis = params.Redis
	p.tokens = params.Tokens

	err = p.db.AutoMigrate(
		Category{},
		Server{},
		ServerCategory{},
	).Error
	return err
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
	local, err := localisation.NewFileSource("assets/translations/serverlist.en.toml", "en")
	if err != nil {
		p.logger.Error("failed to load localisation", zap.Error(err))
	}

	return []interfaces.Localisation{local}
}

func (p *Plugin) Action(event *events.Event) bool {
	if !event.Command() {
		return false
	}

	if event.Fields()[0] != "serverlist" {
		return false
	}

	if len(event.Fields()) >= 2 {
		switch event.Fields()[1] {
		case "category", "categories":

			if len(event.Fields()) >= 3 {
				switch event.Fields()[2] {
				case "create", "add":

					event.Require(func() {

						p.handleCategoryCreate(event)
					},
						permissions.Not(permissions.DiscordChannelDM),
						permissions.BotOwner,
					)
					return true
				}
			}

			p.handleCategoryStatus(event)
			return true

		case "add":

			event.Require(func() {

				p.handleAdd(event)
			},
				permissions.DiscordChannelDM,
			)
			return true

		case "queue":

			if len(event.Fields()) >= 3 {
				if event.Fields()[2] == "refresh" {

					event.Require(func() {

						p.handleQueueRefresh(event)
					},
						permissions.Not(permissions.DiscordChannelDM),
						permissions.BotOwner,
					)
					return true
				}
			}

			event.Require(func() {

				p.handleQueue(event)
			},
				permissions.Not(permissions.DiscordChannelDM),
				permissions.BotOwner,
			)
			return true
		}
	}

	return true
}
