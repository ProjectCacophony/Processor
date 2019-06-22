package customcommands

import (
	"strconv"

	"github.com/jinzhu/gorm"
	"gitlab.com/Cacophony/Processor/plugins/common"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/permissions"
	"gitlab.com/Cacophony/go-kit/regexp"
	"go.uber.org/zap"
)

type Plugin struct {
	logger *zap.Logger
	db     *gorm.DB
}

func (p *Plugin) Name() string {
	return "commands"
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
	return 0
}

func (p *Plugin) Passthrough() bool {
	return false
}

func (p *Plugin) Help() *common.PluginHelp {
	return &common.PluginHelp{
		Name:        p.Name(),
		Description: "customcommands.help.description",
		Commands: []common.Command{{
			Name: "List Commands",
			Params: []common.CommandParam{
				{Name: "list", Type: common.Flag},
				{Name: "public", Type: common.Flag, Optional: true},
			},
		}, {
			Name:        "Add Command",
			Description: "Use the optional parameter 'user', or use command in DM with bot to make a command for just the user.",
			Params: []common.CommandParam{
				{Name: "add", Type: common.Flag},
				{Name: "user", Type: common.Flag, Optional: true},
				{Name: "Command Name", Type: common.Text},
				{Name: "Command Output", Type: common.QuotedText},
			},
		}, {
			Name:        "Edit Command",
			Description: "Use the optional parameter 'user', or use command in DM with bot to edit a command for just the user.",
			Params: []common.CommandParam{
				{Name: "edit", Type: common.Flag},
				{Name: "user", Type: common.Flag, Optional: true},
				{Name: "Command Name", Type: common.Text},
				{Name: "Command Output", Type: common.QuotedText},
			},
		}, {
			Name:        "Delete Command",
			Description: "Use the optional parameter 'user', or use command in DM with bot to delete a command for just the user.",
			Params: []common.CommandParam{
				{Name: "delete", Type: common.Flag},
				{Name: "user", Type: common.Flag, Optional: true},
				{Name: "Command Name", Type: common.Text},
			},
		}, {
			Name: "Toggle Using Commands",
			Description: "Toggle the ability to use the servers custom commands." +
				"\n\t\tAdd 'user' to command to instead toggle the ability for users to use their own custom commands.",
			Params: []common.CommandParam{
				{Name: "toggle", Type: common.Flag},
				{Name: "user", Type: common.Flag, Optional: true},
				// {Name: "channel", Type: common.Channel, Optional: true}, do we need channel specific?
			},
		}, {
			Name:        "Toggle Adding Commands",
			Description: "Enable/disable the ability for everyone, or a specific role, to add new commands.",
			Params: []common.CommandParam{
				{Name: "toggle-permission", Type: common.Flag},
				{Name: "Role ID or Name", Type: common.Text, Optional: true},
			},
		}, {
			Name: "Search for Commands",
			Params: []common.CommandParam{
				{Name: "search", Type: common.Flag},
				{Name: "Command Name", Type: common.Text},
			},
		}, {
			Name:        "Use Random Command",
			Description: "Use a random server command.",
			Params: []common.CommandParam{
				{Name: "random", Type: common.Flag},
				{Name: "user", Type: common.Flag, Optional: true},
			},
		}, {
			Name:        "View Command Info",
			Description: "View information on who added a command, how many times it has been used, and when it was made.",
			Params: []common.CommandParam{
				{Name: "info", Type: common.Flag},
				{Name: "Command Name", Type: common.Text},
			},
		}, {
			Name:            "View Command Info",
			Description:     "View information on who added a command, how many times it has been used, and when it was made.",
			SkipRootCommand: true,
			Params: []common.CommandParam{
				{Name: "Command Name", Type: common.Text},
				{Name: "info", Type: common.Flag},
			},
		}, {
			Name:        "Use User Commad",
			Description: "Generally not needed unless a server command has the same name as your personal command.",
			Params: []common.CommandParam{
				{Name: "user", Type: common.Flag},
				{Name: "Command Name", Type: common.Text},
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
		case "toggle-permissions", "toggle-permission":
			event.Require(func() {
				p.toggleCreatePermission(event)
			}, permissions.DiscordAdministrator)

			return true
		case "toggle":
			event.Require(func() {
				p.toggleUsePermission(event)
			}, permissions.DiscordAdministrator)

			return true
		case "info", "i":
			p.getCommandInfo(event)
			return true
		case "search", "find":
			p.searchCommand(event)
			return true
		case "add":
			p.addCommand(event)
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
