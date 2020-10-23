package autoroles

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/Cacophony/Processor/plugins/common"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/permissions"
	"gitlab.com/Cacophony/go-kit/state"
	"go.uber.org/zap"
)

type Plugin struct {
	logger *zap.Logger
	db     *gorm.DB
	state  *state.State
}

func (p *Plugin) Names() []string {
	return []string{"autorole", "autoroles"}
}

func (p *Plugin) Start(params common.StartParameters) error {
	p.logger = params.Logger
	p.db = params.DB
	p.state = params.State

	err := p.db.AutoMigrate(
		AutoRole{},
	).Error
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
		Names:       p.Names(),
		Description: "autoroles.help.description",
		Commands: []common.Command{{
			Name:        "List Auto role",
			Description: "Lists the currently set auto roles.",
			Params: []common.CommandParam{
				{Name: "list", Type: common.Flag, Optional: true},
			},
		}, {
			Name:        "Add Auto role",
			Description: "Set a role or roles to automatically be assigned to users who join the server.",
			Params: []common.CommandParam{
				{Name: "add", Type: common.Flag},
				{Name: "Role Name", Type: common.QuotedText},
				{Name: "Delay in seconds", Type: common.QuotedText, Optional: true},
			},
		}, {
			Name: "Remove Auto role",
			Params: []common.CommandParam{
				{Name: "remove", Type: common.Flag},
				{Name: "Role Name", Type: common.QuotedText},
			},
		}, {
			Name:        "Auto Role Apply",
			Description: "Applies auto roles to users. Will only apply to users who meet time requirements a if delay exists on the auto role.",
			Params: []common.CommandParam{
				{Name: "apply", Type: common.Flag},
			},
		}},
	}
}

func (p *Plugin) Action(event *events.Event) bool {
	switch event.Type {
	case events.MessageCreateType:
		return p.handleAsCommand(event)
	}

	return false
}

func (p *Plugin) handleAsCommand(event *events.Event) bool {
	if !event.Command() {
		return false
	}

	if !p.isThisModuleCommand(event) {
		return false
	}

	if len(event.Fields()) < 2 {
		p.listAutoRoles(event)
		return true
	}

	if !event.HasOr(permissions.DiscordAdministrator, permissions.DiscordManageRoles) {
		event.Respond("common.missing-role", "roleName", permissions.DiscordManageRoles.Name())
		return true
	}

	switch event.Fields()[1] {
	case "add":
		p.createAutoRole(event)
		return true
	case "remove":
		p.deleteAutoRole(event)
		return true
	case "apply":
		p.applyAutoRole(event)
		return true
	case "list":
		p.listAutoRoles(event)
		return true
	}

	return false
}

func (p Plugin) isThisModuleCommand(event *events.Event) bool {
	for _, name := range p.Names() {
		if event.Fields()[0] == name {
			return true
		}
	}

	return false
}
