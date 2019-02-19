package automod

import (
	"gitlab.com/Cacophony/Processor/plugins/common"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/interfaces"
	"gitlab.com/Cacophony/go-kit/localisation"
)

type Plugin struct {
}

func (p *Plugin) Name() string {
	return "automod"
}

func (p *Plugin) Start(params common.StartParameters) error {
	err := params.DB.AutoMigrate(
		Rule{},
		RuleFilter{},
		RuleAction{},
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
	process := handle(event)

	// if we do not want to further process it, return true to stop further processing
	if !process {
		return true
	}

	// from here on treat as normal command

	if !event.Command() {
		return false
	}

	if event.Fields()[0] != "automod" &&
		event.Fields()[0] != "am" {
		return false
	}

	if len(event.Fields()) < 2 {
		cmdStatus(event)

		return true
	}

	if event.Fields()[1] == "add" {

		cmdAdd(event)
	}

	// TODO: display unknown command message
	return true
}
