package modules

import (
	"fmt"

	"gitlab.com/project-d-collab/dhelpers"
)

// Distributes events to their related modules based on event destination
func CallModules(event dhelpers.EventContainer) {

	for _, module := range moduleList {

		for _, validDest := range module.GetDestinations() {

			for _, targetDest := range event.Destinations {

				if targetDest == validDest {

					// todo: handle panics
					go module.Action(event)
				}
			}
		}
	}
}

// Initializes all plugins
func Init() {
	fmt.Println("Initializing Modules....")

	for _, module := range moduleList {
		module.Init()
	}
}

// Uninitialize all plugins on succesful shutdown
func Uninit() {
	fmt.Println("Uninitializing Modules....")

	for _, module := range moduleList {
		module.Uninit()
	}
}
