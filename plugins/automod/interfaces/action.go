package interfaces

import "gitlab.com/Cacophony/Processor/plugins/automod/models"

type ActionInterface interface {
	Name() string
	Args() int
	NewItem(*models.Env, []string) (ActionItemInterface, error)
	Description() string
	Deprecated() bool
}

type ActionItemInterface interface {
	Do(env *models.Env) (stop bool, err error)
}
