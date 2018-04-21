package modules

import (
	"gitlab.com/project-d-collab/SqsProcessor/modules/color"
	"gitlab.com/project-d-collab/SqsProcessor/modules/ping"
	"gitlab.com/project-d-collab/dhelpers"
)

// Module is an interface for all modules
type Module interface {

	// GetDestinations returns valid destinations for the module
	GetDestinations() []string

	// Init runs at processor startup
	Init()

	// Unit runs at processor shutdown
	Uninit()

	// Action is the main entry point for the module receiving all events
	Action(event dhelpers.EventContainer)
}

var (
	moduleList = []Module{
		&ping.Module{},
		&color.Module{},
	}
)
