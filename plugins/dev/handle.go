package dev

import (
	"gitlab.com/Cacophony/Processor/plugins/common"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/permissions"
	"gitlab.com/Cacophony/go-kit/state"
	"go.uber.org/zap"
)

type Plugin struct {
	state  *state.State
	logger *zap.Logger
}

func (p *Plugin) Name() string {
	return "dev"
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
		Name: p.Name(),
		Hide: true,
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
	}, permissions.DiscordManageServer)
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

		p.handleDevError(event)
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
