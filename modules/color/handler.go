package color

import (
	"fmt"

	"gitlab.com/project-d-collab/dhelpers"
)

type Module struct{}

// Define valid destinations for the module
func (m *Module) GetDestinations() []string {
	return []string{
		"color",
		"colour",
	}
}

// Module entry point when event is triggered
func (m *Module) Action(event dhelpers.EventContainer) {

	fmt.Print("color event triggered")
	switch event.Type {
	case dhelpers.MessageCreateEventType:

		switch event.Args[0] {
		case "color", "colour":
			displayColor(event)
		}
	}
}

// Called on bot startup
func (m *Module) Init() {
	fmt.Println("Initializing color module")
}

// Called on normal bot shutdown
func (m *Module) Uninit() {
	fmt.Println("Uninitializing color module")
}
