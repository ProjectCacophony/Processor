package customcommands

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

func (p *Plugin) Name() string {
	return "commands"
}

func (p *Plugin) Start(params common.StartParameters) error {
	p.logger = params.Logger
	p.db = params.DB

	err := params.DB.AutoMigrate(Entry{}).Error
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
			Name:        "Remove Command",
			Description: "Use the optional parameter 'user', or use command in DM with bot to remove a command for just the user.",
			Params: []common.CommandParam{
				{Name: "remove", Type: common.Flag},
				{Name: "user", Type: common.Flag, Optional: true},
				{Name: "Command Name", Type: common.Text},
			},
		}, {
			Name: "Toggle Using Commands",
			Description: "Toggle the ability to use the servers custom commands in a given channel, or leave out channel to toggle everywhere." +
				"\n\t\tAdd 'user' to command to instead toggle the ability for users to use their own custom commands.",
			Params: []common.CommandParam{
				{Name: "toggle", Type: common.Flag},
				{Name: "user", Type: common.Flag, Optional: true},
				{Name: "channel", Type: common.Channel, Optional: true},
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
			},
		}, {
			Name:        "View Command Info",
			Description: "View information on who added a command, how many times it has been used, and when it was made.",
			Params: []common.CommandParam{
				{Name: "info", Type: common.Flag},
				{Name: "Command Name", Type: common.Text},
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
	if !event.Command() {
		return false
	}

	if event.Fields()[0] != "commands" {
		return p.ranCustomCommand(event)
	}

	if len(event.Fields()) > 1 {
		switch event.Fields()[1] {
		case "add":
			p.addCommand(event)
			return true
		}
	}

	return false
}
