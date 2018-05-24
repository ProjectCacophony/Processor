package gall

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
		"gall",
	}
}

// GetTranslationFiles defines all translation files for the module
func (m *Module) GetTranslationFiles() []string {
	return []string{
		"gall.en.toml",
	}
}

// Action is the module entry point when event is triggered
func (m *Module) Action(ctx context.Context) {
	event := dhelpers.EventFromContext(ctx)

	switch event.Type {
	case dhelpers.MessageCreateEventType:

		for _, destination := range event.Destinations {

			switch destination.Name {
			case "gall":

				if len(event.Args) < 2 {
					return
				}

				switch strings.ToLower(event.Args[1]) {
				case "add": // [p]gall add <board id> [<#target channel or channel id>] [all]
					addBoard(ctx)
					return

				case "list": // [p]gall list
					listFeeds(ctx)
					return

				case "remove", "delete", "rem", "del": // [p]gall remove|delete|rem|delete <board id>
					removeFeed(ctx)
					return

				default: // [p]gall <board id> [all]
					displayBoard(ctx)
					return
				}
			}
		}
	}
}

// Init is called on bot startup
func (m *Module) Init() {
}

// Uninit is called on normal bot shutdown
func (m *Module) Uninit() {
}
