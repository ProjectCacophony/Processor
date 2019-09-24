package roles

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/Cacophony/Processor/plugins/common"
	"gitlab.com/Cacophony/go-kit/events"
	"go.uber.org/zap"
)

type Plugin struct {
	logger *zap.Logger
	db     *gorm.DB
}

func (p *Plugin) Names() []string {
	return []string{"roles", "role"}
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
		Names:       p.Names(),
		Description: "roles.help.description",
		Commands: []common.Command{{
			Name:        "Add Role Category",
			Description: "Add additional role categories for roles to be added too.",
			Params: []common.CommandParam{
				{Name: "add", Type: common.Flag},
				{Name: "category", Type: common.Flag},
				{Name: "Category Name", Type: common.QuotedText},
				{Name: "Category Description", Type: common.QuotedText},
				{Name: "channel", Type: common.Channel},
				{Name: "Limit Count/Pool", Type: common.QuotedText, Optional: true},
			},
		}, {
			Name: "Edit Role Category",
			Params: []common.CommandParam{
				{Name: "edit", Type: common.Flag},
				{Name: "category", Type: common.Flag},
				{Name: "Category Name", Type: common.QuotedText},
				{Name: "Category Description", Type: common.QuotedText},
				{Name: "channel", Type: common.Channel},
				{Name: "Limit Count/Pool", Type: common.QuotedText, Optional: true},
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
				{Name: "Enable/Disable", Type: common.Flag},
				{Name: "category", Type: common.Flag},
				{Name: "Category Name", Type: common.QuotedText},
			},
		}, {
			Name:        "Add Role",
			Description: "Add a role to a given category that users will be able to assign to themselves.",
			Params: []common.CommandParam{
				{Name: "add", Type: common.Flag},
				{Name: "role", Type: common.Flag},
				{Name: "Role Name", Type: common.QuotedText},
				{Name: "Category Name", Type: common.QuotedText, Optional: true},
				{Name: "Print", Type: common.QuotedText, Optional: true},
				{Name: "Aliases...", Type: common.QuotedText, Optional: true},
			},
		}, {
			Name: "Edit Role",
			Params: []common.CommandParam{
				{Name: "edit", Type: common.Flag},
				{Name: "role", Type: common.Flag},
				{Name: "Role Name", Type: common.QuotedText},
				{Name: "Print", Type: common.QuotedText, Optional: true},
				{Name: "Aliases...", Type: common.QuotedText, Optional: true},
			},
		}, {
			Name: "Remove Role",
			Params: []common.CommandParam{
				{Name: "remove", Type: common.Flag},
				{Name: "role", Type: common.Flag},
				{Name: "Role Name", Type: common.QuotedText},
			},
		}, {
			Name:        "Enable/Disable Role",
			Description: "Temporarily disable a role.",
			Params: []common.CommandParam{
				{Name: "Enable/Disable", Type: common.Flag},
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

	// if len(event.Fields()) > 1 {
	// 	switch event.Fields()[1] {
	// 	case "enable":
	// 		p.viewAllPermissions(event)
	// 		return true
	// 	}
	// }

	return false
}
