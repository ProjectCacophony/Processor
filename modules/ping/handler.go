package ping

import (
	"fmt"

	"gitlab.com/project-d-collab/dhelpers"
)

type Module struct{}

// Define valid destinations for the module
func (m *Module) GetDestinations() []string {
	return []string{
		"ping",
	}
}

// Module entry point when event is triggered
func (m *Module) Action(event dhelpers.EventContainer) {

	switch event.Type {
	case dhelpers.MessageCreateEventType:

		switch event.Args[0] {
		case "pong":
			simplePing(event.MessageCreate.ChannelID, event.ReceivedAt)
			break
		case "ping":
			pingInfo(event)
			break
		}
	}
}

// Called on bot startup
func (m *Module) Init() {
	fmt.Println("Initializing ping module")
}

// Called on normal bot shutdown
func (m *Module) Uninit() {
	fmt.Println("Uninitializing ping module")
}
