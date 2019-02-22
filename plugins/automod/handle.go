package automod

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/Cacophony/Processor/plugins/automod/handler"
	"gitlab.com/Cacophony/Processor/plugins/automod/models"
	"gitlab.com/Cacophony/Processor/plugins/common"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/interfaces"
	"gitlab.com/Cacophony/go-kit/localisation"
	"gitlab.com/Cacophony/go-kit/state"
)

type Plugin struct {
	handler *handler.Handler
	state   *state.State
	db      *gorm.DB
}

func (p *Plugin) Name() string {
	return "automod"
}

func (p *Plugin) Start(params common.StartParameters) error {
	err := params.DB.AutoMigrate(
		models.Rule{},
		models.RuleFilter{},
		models.RuleAction{},
	).Error
	if err != nil {
		return err
	}

	p.handler, err = handler.NewHandler(
		params.Logger,
		params.DB,
		params.Redis,
		params.Tokens,
		params.State,
	)
	p.state = params.State
	p.db = params.DB
	return err
}

func (p *Plugin) Stop(params common.StopParameters) error {
	return nil
}

func (p *Plugin) Priority() int {
	return 1000
}

func (p *Plugin) Passthrough() bool {
	return false
}

func (p *Plugin) Localisations() []interfaces.Localisation {
	local, err := localisation.NewFileSource("assets/translations/automod.en.toml", "en")
	if err != nil {
		panic(err) // TODO: handle error
	}

	return []interfaces.Localisation{local}
}

func (p *Plugin) Action(event *events.Event) bool {
	if event.Type != events.MessageCreateType {
		return false
	}

	process := p.handleAsCommand(event)
	if process {
		return true
	}

	// if we do not want to further process it, return true to stop further processing
	return !p.handler.Handle(event)
}

func (p *Plugin) handleAsCommand(event *events.Event) bool {
	if !event.Command() {
		return false
	}

	if event.Fields()[0] != "automod" &&
		event.Fields()[0] != "am" {
		return false
	}
	if len(event.Fields()) < 2 {
		p.cmdStatus(event)

		return true
	}

	switch event.Fields()[1] {

	case "elements", "actions", "filters", "triggers":
		p.cmdElements(event)

		return true
	case "add":

		p.cmdAdd(event)
		return true

	case "remove", "delete":

		p.cmdRemove(event)
		return true
	}

	// TODO: display unknown command message
	return true
}