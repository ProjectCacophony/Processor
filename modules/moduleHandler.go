package modules

import (
	"fmt"
	"runtime"
	"strings"

	"gitlab.com/project-d-collab/dhelpers"
	"gitlab.com/project-d-collab/dhelpers/cache"
)

// CallModules distributes events to their related modules based on event destination
func CallModules(event dhelpers.EventContainer) {

	for _, module := range moduleList {

		for _, validDest := range module.GetDestinations() {

			for _, targetDest := range event.Destinations {

				if targetDest.Name == validDest {

					// send to module
					go func(moduleModule Module, moduleEvent dhelpers.EventContainer) {
						defer func() {
							err := recover()
							if err != nil {
								// error handling
								for _, errorHandlerType := range targetDest.ErrorHandlers {
									switch errorHandlerType {
									case dhelpers.SentryErrorHandler:
										cache.GetLogger().Errorln("handle me via sentry:", err.(error).Error()) // TODO
									case dhelpers.DiscordErrorHandler:
										cache.GetLogger().Errorln("handle me via discord:", err.(error).Error()) // TODO
									}
								}
								// print stacktrace
								buf := make([]byte, 1<<16)
								stackSize := runtime.Stack(buf, false)

								fmt.Println(string(buf[0:stackSize]))
							}
						}()

						moduleModule.Action(moduleEvent)
					}(module, event)
				}
			}
		}
	}
}

// Init initializes all plugins
func Init() {
	cache.GetLogger().Infoln("Initializing Modules....")

	for _, module := range moduleList {
		module.Init()
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
		module.Uninit()
		cache.GetLogger().Infoln("Uninitialized Module for Destinations", "["+strings.Join(module.GetDestinations(), ", ")+"]")
	}
}
