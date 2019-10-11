package dev

import (
	"gitlab.com/Cacophony/Processor/plugins/common"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/interfaces"
	"gitlab.com/Cacophony/go-kit/permissions"
	"gitlab.com/Cacophony/go-kit/state"
	"go.uber.org/zap"
)

type Plugin struct {
	state  *state.State
	logger *zap.Logger
}

func (p *Plugin) Names() []string {
	return []string{"dev"}
}

func (p *Plugin) Start(params common.StartParameters) error {
	p.state = params.State
	p.logger = params.Logger
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
		Names:               p.Names(),
		Description:         "dev.help.description",
		PermissionsRequired: []interfaces.Permission{permissions.BotAdmin},
		Commands: []common.Command{
			{
				Name:        "dev.help.emoji.name",
				Description: "dev.help.emoji.description",
				Params: []common.CommandParam{
					{Name: "emoji", Type: common.Flag},
				},
			},
			{
				Name:        "dev.help.sleep.name",
				Description: "dev.help.sleep.description",
				Params: []common.CommandParam{
					{Name: "sleep", Type: common.Flag},
					{Name: "seconds", Type: common.Text, Optional: true},
				},
			},
			{
				Name:        "dev.help.state.name",
				Description: "dev.help.state.description",
				Params: []common.CommandParam{
					{Name: "state", Type: common.Flag},
					{Name: "user ID", Type: common.User, Optional: true},
				},
			},
			{
				Name:        "dev.help.state-guilds.name",
				Description: "dev.help.state-guilds.description",
				Params: []common.CommandParam{
					{Name: "state", Type: common.Flag},
					{Name: "guilds", Type: common.Flag},
				},
			},
			{
				Name: "dev.help.translate.name",
				Params: []common.CommandParam{
					{Name: "translate", Type: common.Flag},
				},
			},
			{
				Name: "dev.help.error.name",
				Params: []common.CommandParam{
					{Name: "error", Type: common.Flag},
					{Name: "error message", Type: common.Text, Optional: true},
				},
			},
			{
				Name: "dev.help.user-error.name",
				Params: []common.CommandParam{
					{Name: "user-error", Type: common.Flag},
					{Name: "error message", Type: common.Text, Optional: true},
				},
			},
			{
				Name:        "dev.help.permission.name",
				Description: "dev.help.permission.description",
				Params: []common.CommandParam{
					{Name: "permission", Type: common.Flag},
					{Name: "permission code", Type: common.Text, Example: "2048: Send Messages"},
					{Name: "user mention", Type: common.User, Optional: true},
					{Name: "channel mention", Type: common.Channel, Optional: true},
				},
			},
			{
				Name:        "dev.help.questionnaire.name",
				Description: "dev.help.questionnaire.description",
				Params: []common.CommandParam{
					{Name: "questionnaire", Type: common.Flag},
				},
			},
		},
	}
}

func (p *Plugin) Action(event *events.Event) bool {
	if event.Type == events.CacophonyQuestionnaireMatch &&
		event.QuestionnaireMatch.Key == questionnaireKey {
		p.handleDevQuestionnaireMatch(event)
		return true
	}

	if !event.Command() || event.Fields()[0] != "dev" {
		return false
	}

	event.Require(func() {
		p.handleAsCommand(event)
	}, permissions.BotAdmin)
	return true
}

func (p *Plugin) handleAsCommand(event *events.Event) {
	if len(event.Fields()) < 2 {
		event.Respond("dev.no-subcommand")
		return
	}

	switch event.Fields()[1] {
	case "emoji":

		p.handleDevEmoji(event)
		return
	case "sleep":

		p.handleDevSleep(event)
		return
	case "state":
		if len(event.Fields()) > 2 {
			if event.Fields()[2] == "guilds" {
				p.handleDevStateGuilds(event)
				return
			}
		}

		p.handleDevState(event)
		return
	case "translate":

		p.handleDevTranslate(event)
		return
	case "error":

		p.handleDevError(event, false)
		return
	case "user-error":

		p.handleDevError(event, true)
		return
	case "permission":

		p.handleDevPermission(event)
		return
	case "questionnaire":

		p.handleDevQuestionnaire(event)
		return
	}

	event.Respond("dev.no-subcommand")
}
