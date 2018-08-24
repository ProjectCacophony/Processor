package modules

import (
	"context"

	"gitlab.com/Cacophony/Processor/modules/color"
	"gitlab.com/Cacophony/Processor/modules/feed"
	"gitlab.com/Cacophony/Processor/modules/gall"
	"gitlab.com/Cacophony/Processor/modules/lastfm"
	"gitlab.com/Cacophony/Processor/modules/ping"
	"gitlab.com/Cacophony/Processor/modules/stats"
	"gitlab.com/Cacophony/dhelpers"
)

// Module is an interface for all modules
type Module interface {

	// GetDestinations returns valid destinations for the module
	GetDestinations() []string

	// GetTranslationFiles returns all translation files for the module
	GetTranslationFiles() []string

	// Init runs at processor startup
	Init()

	// Unit runs at processor shutdown
	Uninit()

	// Action is the main entry point for the module receiving all events
	Action(ctx context.Context, event dhelpers.EventContainer)
}

var (
	// register modules
	moduleList = []Module{
		&ping.Module{},
		&color.Module{},
		&stats.Module{},
		&lastfm.Module{},
		&gall.Module{},
		&feed.Module{},
	}
)
