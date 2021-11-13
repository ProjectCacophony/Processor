package eventlog

import (
	"strings"

	"github.com/jinzhu/gorm"
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
	return []string{"eventlog"}
}

func (p *Plugin) Start(params common.StartParameters) error {
	p.logger = params.Logger
	p.state = params.State
	p.db = params.DB

	return p.db.AutoMigrate(&Item{}, &ItemOption{}).Error
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
		Description: "eventlog.help.description",
		PermissionsRequired: []interfaces.Permission{
			permissions.DiscordManageServer,
			// permissions.Patron,
		},
		Commands: []common.Command{
			{
				Name:        "eventlog.help.status.name",
				Description: "eventlog.help.status.description",
			},
			{
				Name:        "eventlog.help.history.name",
				Description: "eventlog.help.history.description",
				Params: []common.CommandParam{
					{Name: "history", Type: common.Flag},
					{Name: "user or user id", Type: common.User},
				},
			},
			{
				Name:        "eventlog.help.enable.name",
				Description: "eventlog.help.enable.description",
				Params: []common.CommandParam{
					{Name: "enable", Type: common.Flag},
					{Name: "log channel", Type: common.Channel, Optional: true},
				},
				PermissionsRequired: []interfaces.Permission{permissions.DiscordAdministrator},
			},
			{
				Name:        "eventlog.help.disable.name",
				Description: "eventlog.help.disable.description",
				Params: []common.CommandParam{
					{Name: "disable", Type: common.Flag},
				},
				PermissionsRequired: []interfaces.Permission{permissions.DiscordAdministrator},
			},
		},
	}
}

func (p *Plugin) Action(event *events.Event) bool {
	switch event.Type {
	case events.CacophonyEventlogUpdate:
		p.handleEventlogUpdate(event)
		return true
	case events.MessageReactionAddType:
		event.Require(func() {
			p.handleReactionAdd(event)
		},
			permissions.DiscordManageServer,
			permissions.Not(
				permissions.DiscordChannelDM,
			),
		)
	case events.CacophonyQuestionnaireMatch:
		if event.QuestionnaireMatch.Key == questionnaireEditReasonKey {
			p.handleQuestionnaireEditReason(event)
			return true
		}
	}

	p.handleModEvent(event)

	if !event.Command() {
		return false
	}

	if event.Fields()[0] != "eventlog" {
		return false
	}

	event.Require(
		func() {
			p.handleCommand(event)
		},
		permissions.DiscordManageServer,
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

			event.Require(func() {
				p.handleCmdEnable(event)
			}, permissions.DiscordAdministrator) // , permissions.Patron)
			return
		case "disable":

			event.Require(func() {
				p.handleCmdDisable(event)
			}, permissions.DiscordAdministrator) // , permissions.Patron)
			return

		case "history":
			p.handleCmdHistory(event)
			return
		}
	}

	p.handleCmdStatus(event)
}
