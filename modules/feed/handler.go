package feed

import (
	"strings"

	"gitlab.com/project-d-collab/dhelpers"
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
func (m *Module) Action(event dhelpers.EventContainer) {

	switch event.Type {
	case dhelpers.MessageCreateEventType:

		for _, destination := range event.Destinations {

			switch destination.Name {
			case "feed":

				if len(event.Args) < 2 {
					return
				}

				switch strings.ToLower(event.Args[1]) {
				case "add": // [p]feed add <board id> [<#target channel or channel id>]
					//addBoard(event)
					return

				case "list": // [p]feed list
					//listFeeds(event)
					return

				case "remove", "delete", "rem", "del": // [p]feed remove|delete|rem|delete <board id>
					//removeFeed(event)
					return

				default: // [p]feed <feed url>
					displayFeed(event)
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
