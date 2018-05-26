package ping

import (
	"strings"

	"context"

	"gitlab.com/Cacophony/dhelpers"
)

// Module is a struct for the module
type Module struct{}

// GetDestinations defines valid destinations for the module
func (m *Module) GetDestinations() []string {
	return []string{
		"ping",
	}
}

// GetTranslationFiles defines all translation files for the module
func (m *Module) GetTranslationFiles() []string {
	return []string{}
}

// Action is the module entry point when event is triggered
func (m *Module) Action(ctx context.Context, event dhelpers.EventContainer) {
	switch event.Type {
	case dhelpers.MessageCreateEventType:

		switch strings.ToLower(event.Args[0]) {
		case "pong":
			simplePing(ctx, event, event.ReceivedAt)
		case "ping":
			pingInfo(ctx, event)
		}
	}
}

// Init is called on bot startup
func (m *Module) Init() {
}

// Uninit is called on normal bot shutdown
func (m *Module) Uninit() {
}
