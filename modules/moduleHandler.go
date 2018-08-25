package modules

import (
	"strings"

	"context"

	"fmt"

	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
	"gitlab.com/Cacophony/dhelpers"
	"gitlab.com/Cacophony/dhelpers/cache"
)

// CallModules distributes events to their related modules based on event destination
func CallModules(event dhelpers.EventContainer) {

	for _, module := range moduleList {

		for _, validDest := range module.GetDestinations() {

			for _, targetDest := range event.Destinations {

				if targetDest.Name == validDest {

					// send to module
					go func(destination string, moduleModule Module, moduleEvent dhelpers.EventContainer) {
						defer func() {
							err := recover()
							if err != nil {
								if _, ok := err.(error); !ok {
									err = fmt.Errorf("%+v", err)
								}
								// handle errors
								dhelpers.HandleErrWith("Processor", err.(error), &moduleEvent, targetDest.ErrorHandlers...)
							}
						}()

						// start span from background
						var span opentracing.Span
						ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
						span, ctx = opentracing.StartSpanFromContext(ctx, destination)
						// add fields
						span.LogFields(
							log.String("event", event.Key),
						)

						defer span.Finish()
						defer cancel()

						// start action
						moduleModule.Action(ctx, event)
					}(targetDest.Name, module, event)
				}
			}
		}
	}
}

// Init initializes all plugins
func Init() {
	cache.GetLogger().Infoln("Initializing Modules....")

	for _, module := range moduleList {
		// initialise module
		module.Init()
		// load translation files for module
		for _, translationFileName := range module.GetTranslationFiles() {
			_, err := cache.GetLocalizationBundle().LoadMessageFile("./translations/" + translationFileName)
			if err != nil {
				panic(err)
			}
			cache.GetLogger().Infoln("Loaded " + translationFileName)
		}
		cache.GetLogger().Infoln("Initialized Module for Destinations", "["+strings.Join(module.GetDestinations(), ", ")+"]")
	}
}

// Uninit uninitialize all plugins on succesfull shutdown
func Uninit() {
	cache.GetLogger().Infoln("Uninitializing Modules....")

	for _, module := range moduleList {
		// uninitialise module
		module.Uninit()
		cache.GetLogger().Infoln("Uninitialized Module for Destinations", "["+strings.Join(module.GetDestinations(), ", ")+"]")
	}
}
