package roles

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
	return []string{"roles", "role"}
}

func (p *Plugin) Start(params common.StartParameters) error {
	p.logger = params.Logger
	p.db = params.DB
	p.state = params.State

	err := p.db.AutoMigrate(
		Category{},
		Role{},
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
		Description: "roles.help.description",
		Commands: []common.Command{{
			Name:        "Add Role Category",
			Description: "Add additional role categories for roles to be added too.",
			Params: []common.CommandParam{
				{Name: "add", Type: common.Flag},
				{Name: "category", Type: common.Flag},
				{Name: "Category Name", Type: common.QuotedText},
				{Name: "channel", Type: common.Channel, Optional: true},
				{Name: "Limit Count", Type: common.QuotedText, Optional: true},
			},
		}, {
			Name: "Edit Role Category",
			Params: []common.CommandParam{
				{Name: "edit", Type: common.Flag},
				{Name: "category", Type: common.Flag},
				{Name: "Category Name", Type: common.QuotedText},
				{Name: "New Category Name", Type: common.QuotedText},
				{Name: "channel", Type: common.Channel, Optional: true},
				{Name: "Limit Count", Type: common.QuotedText, Optional: true},
			},
		}, {
			Name: "Remove Role Category",
			Params: []common.CommandParam{
				{Name: "remove", Type: common.Flag},
				{Name: "category", Type: common.Flag},
				{Name: "Category Name", Type: common.QuotedText},
			},
		}, {
			Name:        "Enable/Disable Role Category",
			Description: "Temporarily disable a role category.",
			Params: []common.CommandParam{
				{Name: "enable/disable", Type: common.Flag},
				{Name: "category", Type: common.Flag},
				{Name: "Category Name", Type: common.QuotedText},
			},
		}, {
			Name:        "Show/Hide Role Category",
			Description: "Toggle whether or not the category and its roles will show in the role channel text. Unlike enabling/disabling, roles from hidden categories can still be assigned.",
			Params: []common.CommandParam{
				{Name: "show/hide", Type: common.Flag},
				{Name: "category", Type: common.Flag},
				{Name: "Category Name", Type: common.QuotedText},
			},
		}, {
			Name:        "Add Role",
			Description: "Add a role to a given category that users will be able to assign to themselves.",
			Params: []common.CommandParam{
				{Name: "add", Type: common.Flag},
				{Name: "role", Type: common.Flag},
				{Name: "Role Name/ID", Type: common.QuotedText},
				{Name: "Print", Type: common.QuotedText, Optional: true},
				{Name: "Alias, Alias...", Type: common.QuotedText, Optional: true},
				{Name: "Category Name", Type: common.QuotedText, Optional: true},
			},
		}, {
			Name: "Edit Role",
			Params: []common.CommandParam{
				{Name: "edit", Type: common.Flag},
				{Name: "role", Type: common.Flag},
				{Name: "Role Name/ID", Type: common.QuotedText},
				{Name: "Print", Type: common.QuotedText, Optional: true},
				{Name: "Alias, Alias...", Type: common.QuotedText, Optional: true},
				{Name: "Category Name", Type: common.QuotedText, Optional: true},
			},
		}, {
			Name: "Remove Role",
			Params: []common.CommandParam{
				{Name: "remove", Type: common.Flag},
				{Name: "role", Type: common.Flag},
				{Name: "Role Name", Type: common.QuotedText},
			},
		}, {
			Name:        "Role Info",
			Description: "view current role categories and roles confirgured on the server.",
			Params: []common.CommandParam{
				{Name: "info", Type: common.Flag},
			},
		}, {
			Name:        "Add Auto role",
			Description: "Set a role or roles to automatically be assigned to users who join the server.",
			Params: []common.CommandParam{
				{Name: "auto", Type: common.Flag},
				{Name: "add", Type: common.QuotedText},
				{Name: "Role Name", Type: common.QuotedText},
				{Name: "Delay in minutes/hours/days", Type: common.QuotedText, Optional: true},
			},
		}, {
			Name: "Remove Auto role",
			Params: []common.CommandParam{
				{Name: "auto", Type: common.Flag},
				{Name: "remove", Type: common.QuotedText},
				{Name: "Role Name", Type: common.QuotedText},
			},
		}, {
			Name:        "List Auto role",
			Description: "Lists our the currently set auto roles.",
			Params: []common.CommandParam{
				{Name: "auto", Type: common.Flag},
				{Name: "list", Type: common.QuotedText},
			},
		}, {
			Name:        "Auto Role Apply",
			Description: "Applies auto roles to users. Will only apply to users who meet time requirements a delay exists on the auto role.",
			Params: []common.CommandParam{
				{Name: "auto", Type: common.Flag},
				{Name: "list", Type: common.QuotedText},
			},
		}, {
			Name:        "Set Role Channel",
			Description: "Set the default channel users will go to to add and remove their roles.",
			Params: []common.CommandParam{
				{Name: "channel", Type: common.Flag},
				{Name: "channel", Type: common.Channel},
			},
		}},
	}
}

func (p *Plugin) Action(event *events.Event) bool {
	if !event.Command() {
		return false
	}

	if event.Fields()[0] != "roles" &&
		event.Fields()[0] != "role" {
		return false
	}

	if len(event.Fields()) > 1 {
		switch event.Fields()[1] {
		case "show", "hide":
			p.toggleCategoryVisibility(event)
			return true
		case "enable", "disable":
			p.toggleCategory(event)
			return true
		case "channel":
			p.setRoleChannel(event)
			return true
		case "info":
			p.displayRoleInfo(event)
			return true
		case "add":
			if len(event.Fields()) < 3 {
				return true
			}

			if !event.HasOr(permissions.DiscordAdministrator, permissions.DiscordManageRoles) {
				event.Respond("common.missing-role", "roleName", permissions.DiscordManageRoles.Name())
				return true
			}

			switch event.Fields()[2] {
			case "category":

				p.createCategory(event)
				return true
			case "role":
				p.createRole(event)
				return true
			}
		case "edit", "update":
			if len(event.Fields()) < 3 {
				return true
			}

			if !event.HasOr(permissions.DiscordAdministrator, permissions.DiscordManageRoles) {
				event.Respond("common.missing-role", "roleName", permissions.DiscordManageRoles.Name())
				return true
			}

			switch event.Fields()[2] {
			case "category":

				p.updateCategory(event)
				return true
			case "role":
				p.updateRole(event)
				return true
			}
			return true

		case "delete", "remove":
			if len(event.Fields()) < 3 {
				return true
			}

			if !event.HasOr(permissions.DiscordAdministrator, permissions.DiscordManageRoles) {
				event.Respond("common.missing-role", "roleName", permissions.DiscordManageRoles.Name())
				return true
			}

			switch event.Fields()[2] {
			case "category":

				p.deleteCategory(event)
				return true
			case "role":
				p.deleteRole(event)
				return true
			}
			return true

		default:
		}
	}

	return false
}
