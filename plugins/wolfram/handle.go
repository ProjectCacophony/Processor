package wolfram

import (
	"errors"

	"github.com/Krognol/go-wolfram"
	"github.com/jinzhu/gorm"
	"github.com/kelseyhightower/envconfig"
	"gitlab.com/Cacophony/Processor/plugins/common"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/state"
	"go.uber.org/zap"
)

type Plugin struct {
	logger        *zap.Logger
	db            *gorm.DB
	state         *state.State
	wolframClient *wolfram.Client
}

func (p *Plugin) Names() []string {
	return []string{"wolfram", "w"}
}

func (p *Plugin) Start(params common.StartParameters) error {
	p.logger = params.Logger
	p.db = params.DB
	p.state = params.State

	var config Config
	err := envconfig.Process("", &config)
	if err != nil || config.WolframAPPID == "" {
		return errors.New("Unable to process or find wolfram alpha environment variables")
	}

	p.wolframClient = &wolfram.Client{AppID: config.WolframAPPID}

	return nil
}

func (p *Plugin) Stop(params common.StopParameters) error {
	return nil
}

func (p *Plugin) Priority() int {
	return 100
}

func (p *Plugin) Passthrough() bool {
	return false
}

func (p *Plugin) Help() *common.PluginHelp {
	return &common.PluginHelp{
		Names:       p.Names(),
		Description: "wolfram.description",
		Commands: []common.Command{{
			Name:        "Ask a question",
			Description: "",
			Params: []common.CommandParam{
				{Name: "question...", Type: common.QuotedText},
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

	if event.Fields()[0] != "wolfram" &&
		event.Fields()[0] != "w" {
		return false
	}

	event.Typing()
	p.askWolfram(event)
	return true
}
