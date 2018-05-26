package feed

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
		"feed",
	}
}

// GetTranslationFiles defines all translation files for the module
func (m *Module) GetTranslationFiles() []string {
	return []string{
		"feed.en.toml",
	}
}

// Action is the module entry point when event is triggered
func (m *Module) Action(ctx context.Context, event dhelpers.EventContainer) {
	switch event.Type {
	case dhelpers.MessageCreateEventType:

		for _, destination := range event.Destinations {

			switch destination.Name {
			case "feed":

				if len(event.Args) < 2 {
					return
				}

				switch strings.ToLower(event.Args[1]) {
				case "add": // [p]feed add <feed url> [<#target channel or channel id>]
					addFeed(ctx, event)
					return

				case "list": // [p]feed list
					listFeeds(ctx, event)
					return

				case "remove", "delete", "rem", "del": // [p]feed remove|delete|rem|delete <board id>
					removeFeed(ctx, event)
					return

				default: // [p]feed <feed url>
					displayFeed(ctx, event)
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
