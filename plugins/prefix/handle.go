package prefix

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/Cacophony/Processor/plugins/common"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/interfaces"
	"gitlab.com/Cacophony/go-kit/permissions"
	"go.uber.org/zap"
)

type Plugin struct {
	logger *zap.Logger
	db     *gorm.DB
}

func (p *Plugin) Name() string {
	return "prefix"
}

func (p *Plugin) Start(params common.StartParameters) error {
	p.logger = params.Logger
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

func (p *Plugin) Help() *common.PluginHelp {
	return &common.PluginHelp{
		Name:        p.Name(),
		Description: "prefix.help.description",
		Commands: []common.Command{{
			Name: "View Prefix",
		}, {
			Name:                "Set Prefix",
			PermissionsRequired: []interfaces.Permission{permissions.DiscordAdministrator},
			Params: []common.CommandParam{
				{Name: "set", Type: common.Flag},
				{Name: "new prefix", Type: common.Text},
			},
		}},
	}
}

func (p *Plugin) Action(event *events.Event) bool {
	if !event.Command() {
		return false
	}

	if event.Fields()[0] == "prefix" {

		if len(event.Fields()) == 1 {
			handleGetPrefix(event)
			return true
		}

		if event.Fields()[1] == "set" {

			event.Require(func() {

				handleSetPrefix(event, p.db)
			}, permissions.DiscordAdministrator)

			return true
		}
	}

	return false
}
