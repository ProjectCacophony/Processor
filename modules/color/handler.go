package color

import (
	"gitlab.com/project-d-collab/dhelpers"
)

// Module is a struct for the module
type Module struct{}

// GetDestinations defines valid destinations for the module
func (m *Module) GetDestinations() []string {
	return []string{
		"color",
		"colour",
	}
}

// Action is the module entry point when event is triggered
func (m *Module) Action(event dhelpers.EventContainer) {
	switch event.Type {
	case dhelpers.MessageCreateEventType:

		switch event.Args[0] {
		case "color", "colour":
			displayColor(event)
		}
	}
}

// Init is called on bot startup
func (m *Module) Init() {
}

// Uninit is called on normal bot shutdown
func (m *Module) Uninit() {
}
