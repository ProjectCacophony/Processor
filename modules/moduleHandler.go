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

				if targetDest.Name == validDest {

					// send to module
					go func(moduleModule Module, moduleDequeudAt time.Time, moduleEvent dhelpers.EventContainer) {
						defer func() {
							err := recover()
							if err != nil {
								for _, errorHandlerType := range targetDest.ErrorHandlers {
									switch errorHandlerType {
									case dhelpers.SentryErrorHandler:
										fmt.Printf("handle me via sentry: %+v\n", err) // TODO
									case dhelpers.DiscordErrorHandler:
										fmt.Printf("handle me via discord: %+v\n", err) // TODO
									}
								}
							}
						}()

						moduleModule.Action(moduleDequeudAt, moduleEvent)
					}(module, dequeuedAt, event)
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
