package modules

import (
	"fmt"

	"strings"

	"time"

	"gitlab.com/project-d-collab/dhelpers"
)

// CallModules distributes events to their related modules based on event destination
func CallModules(dequeuedAt time.Time, event dhelpers.EventContainer) {

	for _, module := range moduleList {

		for _, validDest := range module.GetDestinations() {

			for _, targetDest := range event.Destinations {

				if targetDest == validDest {

					// todo: handle panics
					go module.Action(dequeuedAt, event)
				}
			}
		}
	}
}

// Init initializes all plugins
func Init() {
	fmt.Println("Initializing Modules....")

	for _, module := range moduleList {
		module.Init()
		fmt.Println("Initialized Module for Destinations", "["+strings.Join(module.GetDestinations(), ", ")+"]")
	}
}

// Uninit uninitialize all plugins on succesfull shutdown
func Uninit() {
	fmt.Println("Uninitializing Modules....")

	for _, module := range moduleList {
		module.Uninit()
		fmt.Println("Uninitialized Module for Destinations", "["+strings.Join(module.GetDestinations(), ", ")+"]")
	}
}
