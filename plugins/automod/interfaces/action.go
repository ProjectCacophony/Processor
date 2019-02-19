package interfaces

import "gitlab.com/Cacophony/Processor/plugins/automod/models"

type ActionInterface interface {
	Name() string
	NewItem(string) (ActionItemInterface, error)
}

type ActionItemInterface interface {
	Do(env *models.Env)
}
