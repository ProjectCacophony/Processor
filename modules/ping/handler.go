package ping

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
		"ping",
	}
}

// Action is the module entry point when event is triggered
func (m *Module) Action(dequeuedAt time.Time, event dhelpers.EventContainer) {

	switch event.Type {
	case dhelpers.MessageCreateEventType:

		switch event.Args[0] {
		case "pong":
			simplePing(event.MessageCreate.ChannelID, event.ReceivedAt)
		case "ping":
			pingInfo(dequeuedAt, event)
		}
	}
}

// Init is called on bot startup
func (m *Module) Init() {
	fmt.Println("Initializing ping module")
}

// Uninit is called on normal bot shutdown
func (m *Module) Uninit() {
	fmt.Println("Uninitializing ping module")
}
