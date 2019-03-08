package chatlog

import (
	"strings"

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
	state  *state.State
	db     *gorm.DB
}

func (p *Plugin) Name() string {
	return "chatlog"
}

func (p *Plugin) Start(params common.StartParameters) error {
	p.logger = params.Logger
	p.state = params.State
	p.db = params.DB
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
	local, err := localisation.NewFileSource("assets/translations/chatlog.en.toml", "en")
	if err != nil {
		p.logger.Error("failed to load localisation", zap.Error(err))
	}

	return []interfaces.Localisation{local}
}

func (p *Plugin) Action(event *events.Event) bool {
	if !event.Command() {
		return false
	}

	if event.Fields()[0] != "chatlog" {
		return false
	}

	// has to be Server Admin, and Patron/Staff
	event.Require(
		func() {
			p.handleCommand(event)
		},
		permissions.DiscordAdministrator,
		permissions.Or(
			// TODO: don't hardcode
			// sekl's dev cord / Cacophony Whitelist Role
			permissions.NewDiscordRole(
				p.state, "208673735580844032", "549951245738180688",
			),
			// Cacophony / Cacophony Team
			permissions.NewDiscordRole(
				p.state, "435420687906111498", "440519691904090113",
			),
		),
		permissions.Not(
			permissions.DiscordChannelDM,
		),
	)

	return true
}

func (p *Plugin) handleCommand(event *events.Event) {

	if len(event.Fields()) >= 2 {
		switch strings.ToLower(event.Fields()[1]) {

		case "enable":

			p.handleEnable(event)
		case "disable":

			p.handleDisable(event)
		}
	}

	p.handleStatus(event)
}
