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
	logger     *zap.Logger
	db         *gorm.DB
	state      *state.State
	redis      *redis.Client
	tokens     map[string]string
	staffRoles permissions.PermissionInterface
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

	p.staffRoles = permissions.Or(
		// sekl's dev cord / Admin
		permissions.NewDiscordRole(p.state, "208673735580844032", "250710478068645890"),
		// Test / Staff
		permissions.NewDiscordRole(p.state, "561619599129444390", "561619665197989893"),
	)

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

func (p *Plugin) Help() *common.PluginHelp {
	return &common.PluginHelp{
		Name:        p.Name(),
		Description: "Uhh.... sekl what is this even doing :nayoungconfused:",
	}
}

func (p *Plugin) Action(event *events.Event) bool {
	switch event.Type {

	case events.MessageReactionAddType:
		return p.handleQueueReaction(event)

	case events.CacophonyServerlistServerExpire:
		p.handleExpired(event)
		return true

	}

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
						p.staffRoles,
					)
					return true
				}
			}

			p.handleCategoryStatus(event)
			return true

		case "add":

			p.handleAdd(event)
			return true

		case "queue":

			if len(event.Fields()) >= 3 {
				if event.Fields()[2] == "refresh" {

					event.Require(func() {

						p.handleQueueRefresh(event)
					},
						permissions.Not(permissions.DiscordChannelDM),
						p.staffRoles,
					)
					return true
				}
			}

			event.Require(func() {

				p.handleQueue(event)
			},
				permissions.Not(permissions.DiscordChannelDM),
				p.staffRoles,
			)
			return true

		case "list":

			if len(event.Fields()) >= 3 {

				switch event.Fields()[2] {

				case "refresh":

					event.Require(func() {

						p.handleListRefresh(event)
					},
						permissions.Not(permissions.DiscordChannelDM),
						p.staffRoles,
					)
					return true

				case "clear-cache":

					event.Require(func() {

						p.handleListClearCache(event)
					},
						permissions.Not(permissions.DiscordChannelDM),
						p.staffRoles,
					)
					return true

				}
			}

		case "reject":

			event.Require(func() {

				p.handleQueueReject(event)
			},
				permissions.Not(permissions.DiscordChannelDM),
				p.staffRoles,
			)
			return true

		case "remove":

			p.handleRemove(event)
			return true

		case "hide", "unhide":

			p.handleHide(event)
			return true

		case "log":

			event.Require(func() {

				p.handleLog(event)
			},
				permissions.Not(permissions.DiscordChannelDM),
				p.staffRoles,
			)
			return true

		case "edit":

			p.handleEdit(event)
			return true

		}
	}

	p.handleStatus(event)

	return true
}
