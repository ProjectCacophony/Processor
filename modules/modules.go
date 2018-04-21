package modules

import (
	"gitlab.com/project-d-collab/SqsProcessor/modules/color"
	"gitlab.com/project-d-collab/SqsProcessor/modules/ping"
	"gitlab.com/project-d-collab/dhelpers"
)

type Module interface {

	// returns valid destinations for the module
	GetDestinations() []string

	// Bot startup/shutdown
	Init()
	Uninit()

	// Main entry point for module
	Action(event dhelpers.EventContainer)
}

var (
	moduleList = []Module{
		&ping.Module{},
		&color.Module{},
	}
)
