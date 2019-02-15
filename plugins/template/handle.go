package template

import (
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/interfaces"
	"gitlab.com/Cacophony/go-kit/localisation"
)

type Plugin struct{}

func (p *Plugin) Name() string {
	return "<plugin name>" // TODO: change me
}

func (p *Plugin) Start() error {
	return nil
}

func (p *Plugin) Stop() error {
	return nil
}

func (p *Plugin) Priority() int {
	return 0
}

func (p *Plugin) Passthrough() bool {
	return false
}

func (p *Plugin) Localisations() []interfaces.Localisation {
	local, err := localisation.NewFileSource("assets/translations/<file name>.en.toml", "en") // TODO: change me
	if err != nil {
		panic(err) // TODO: handle error
	}

	return []interfaces.Localisation{local}
}

func (p *Plugin) Action(event *events.Event) bool {
	if !event.Command() {
		return false
	}

	switch event.Fields()[0] {
	case "foo":

		handleFoo(event)
		return true
	case "bar":

		handleBar(event)
		return true
	}

	return false
}
