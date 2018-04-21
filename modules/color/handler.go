package color

import (
	"fmt"

	"time"

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
func (m *Module) Action(_ time.Time, event dhelpers.EventContainer) {

	fmt.Println("color event triggered")
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
	fmt.Println("Initializing color module")
}

// Uninit is called on normal bot shutdown
func (m *Module) Uninit() {
	fmt.Println("Uninitializing color module")
}
