package customcommands

import (
	"strconv"

	"github.com/jinzhu/gorm"
	"gitlab.com/Cacophony/Processor/plugins/common"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/interfaces"
	"gitlab.com/Cacophony/go-kit/permissions"
	"gitlab.com/Cacophony/go-kit/regexp"
	"go.uber.org/zap"
)

type Plugin struct {
	logger *zap.Logger
	db     *gorm.DB
}

func (p *Plugin) Names() []string {
	return []string{"commands"}
}

func (p *Plugin) Start(params common.StartParameters) error {
	p.logger = params.Logger
	p.db = params.DB

	err := params.DB.AutoMigrate(CustomCommand{}).Error
	if err != nil {
		return err
	}

	return nil
}

func (p *Plugin) Stop(params common.StopParameters) error {
	return nil
}

func (p *Plugin) Priority() int {
	return -1000
}

func (p *Plugin) Passthrough() bool {
	return false
}

func (p *Plugin) Help() *common.PluginHelp {
	return &common.PluginHelp{
		Names:       p.Names(),
		Description: "customcommands.help.description",
		Commands: []common.Command{{
			Name: "customcommands.help.list.name",
			Params: []common.CommandParam{
				{Name: "list", Type: common.Flag},
				{Name: "public", Type: common.Flag, Optional: true},
			},
		}, {
			Name:        "customcommands.help.add.name",
			Description: "customcommands.help.add.description",
			Params: []common.CommandParam{
				{Name: "add", Type: common.Flag},
				{Name: "user", Type: common.Flag, Optional: true},
				{Name: "Command Names", Type: common.Text},
				{Name: "Command Output", Type: common.QuotedText},
			},
		}, {
			Name:        "customcommands.help.add-alias.name",
			Description: "customcommands.help.add-alias.description",
			Params: []common.CommandParam{
				{Name: "add-alias", Type: common.Flag},
				{Name: "user", Type: common.Flag, Optional: true},
				{Name: "Command Alias", Type: common.Text},
				{Name: "Actual Command", Type: common.QuotedText},
			},
			PermissionsRequired: []interfaces.Permission{permissions.DiscordManageServer},
		}, {
			Name:        "customcommands.help.edit.name",
			Description: "customcommands.help.edit.description",
			Params: []common.CommandParam{
				{Name: "edit", Type: common.Flag},
				{Name: "user", Type: common.Flag, Optional: true},
				{Name: "Command Names", Type: common.Text},
				{Name: "Command Output", Type: common.QuotedText},
			},
		}, {
			Name:        "customcommands.help.delete.name",
			Description: "customcommands.help.delete.description",
			Params: []common.CommandParam{
				{Name: "delete", Type: common.Flag},
				{Name: "user", Type: common.Flag, Optional: true},
				{Name: "Command Names", Type: common.Text},
			},
		}, {
			Name:        "customcommands.help.viewpermission.name",
			Description: "customcommands.help.viewpermission.description",
			Params: []common.CommandParam{
				{Name: "permissions", Type: common.Flag},
			},
		}, {
			Name:        "customcommands.help.usepermission.name",
			Description: "customcommands.help.usepermission.description",
			Params: []common.CommandParam{
				{Name: "enable/disable", Type: common.Flag},
				{Name: "user", Type: common.Flag, Optional: true},
				// {Names: "channel", Type: common.Channel, Optional: true}, do we need channel specific?
			},
		}, {
			Name:        "customcommands.help.addpermission.name",
			Description: "customcommands.help.addpermission.description",
			Params: []common.CommandParam{
				{Name: "enable/disable", Type: common.Flag},
				{Name: "edit", Type: common.Flag},
				{Name: "Role ID or Names", Type: common.Text, Optional: true},
			},
		}, {
			Name: "customcommands.help.search.name",
			Params: []common.CommandParam{
				{Name: "search", Type: common.Flag},
				{Name: "Command Names", Type: common.Text},
			},
		}, {
			Name:        "customcommands.help.random.name",
			Description: "customcommands.help.random.description",
			Params: []common.CommandParam{
				{Name: "random", Type: common.Flag},
				{Name: "user", Type: common.Flag, Optional: true},
			},
		}, {
			Name:            "customcommands.help.info.name",
			Description:     "customcommands.help.info.description",
			SkipRootCommand: true,
			Params: []common.CommandParam{
				{Name: "Command Names", Type: common.Text},
				{Name: "info", Type: common.Flag},
			},
		}, {
			Name:        "customcommands.help.user.name",
			Description: "customcommands.help.user.description",
			Params: []common.CommandParam{
				{Name: "user", Type: common.Flag},
				{Name: "Command Names", Type: common.Text},
			},
		}},
	}
}

func (p *Plugin) Action(event *events.Event) bool {
	if event.Type == events.CacophonyQuestionnaireMatch &&
		(event.QuestionnaireMatch.Key == editQuestionnaireKey ||
			event.QuestionnaireMatch.Key == deleteQuestionnaireKey) {
		p.handleQuestionnaire(event)
		return true
	}

	if !event.Command() {
		return false
	}

	if event.Fields()[0] != "commands" {
		return p.runCustomCommand(event)
	}

	if len(event.Fields()) > 1 {
		switch event.Fields()[1] {
		case "enable":
			if len(event.Fields()) > 2 && event.Fields()[2] == "edit" {
				event.Require(func() {
					p.toggleCreatePermission(event, true)
				}, permissions.DiscordAdministrator)
				return true
			}

			event.Require(func() {
				p.toggleUsePermission(event, false)
			}, permissions.DiscordAdministrator)

			return true
		case "disable":
			if len(event.Fields()) > 2 && event.Fields()[2] == "edit" {
				event.Require(func() {
					p.toggleCreatePermission(event, false)
				}, permissions.DiscordAdministrator)
				return true
			}

			event.Require(func() {
				p.toggleUsePermission(event, true)
			}, permissions.DiscordAdministrator)

			return true
		case "info", "i":
			p.getCommandInfo(event)
			return true
		case "search", "find":
			p.searchCommand(event)
			return true
		case "add":
			p.addCommand(event, customCommandTypeContent)
			return true
		case "add-alias":
			p.addCommand(event, customCommandTypeCommand)
			return true
		case "edit":
			p.editCommand(event)
			return true
		case "list":
			p.listCommands(event)
			return true
		case "delete", "del", "remove":
			p.deleteCommand(event)
			return true

		case "rand", "random":
			p.runRandomCommand(event)
			return true
		case "permissions":
			p.viewAllPermissions(event)
			return true
		}
	}

	return false
}

func (p *Plugin) handleQuestionnaire(event *events.Event) {
	if event.MessageCreate == nil || event.MessageCreate.Content == "" {
		return
	}

	var handled bool
	if enteredNum, err := strconv.Atoi(regexp.ContainsNumber.FindString(event.MessageCreate.Content)); err == nil {
		switch event.QuestionnaireMatch.Key {
		case editQuestionnaireKey:
			handled = p.handleEditResponse(event, enteredNum)
		case deleteQuestionnaireKey:
			handled = p.handleDeleteResponse(event, enteredNum)
		}
		if !handled {
			event.Questionnaire().Redo(event)
		}
	} else {
		event.Questionnaire().Redo(event)
	}
}
