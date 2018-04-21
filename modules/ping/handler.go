package ping

import (
	"gitlab.com/project-d-collab/dhelpers"
	"gitlab.com/project-d-collab/dhelpers/cache"
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
func (m *Module) Action(event dhelpers.EventContainer) {

	cache.GetLogger().Infoln("Ping module triggered at: ", event.ReceivedAt)
	switch event.Type {
	case dhelpers.MessageCreateEventType:

		switch event.Args[0] {
		case "pong":
			simplePing(event.MessageCreate.ChannelID, event.ReceivedAt)
		case "ping":
			pingInfo(event)
		}
	}
}

// Init is called on bot startup
func (m *Module) Init() {
}

// Uninit is called on normal bot shutdown
func (m *Module) Uninit() {
}
