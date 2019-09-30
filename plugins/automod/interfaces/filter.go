package interfaces

import (
	"gitlab.com/Cacophony/Processor/plugins/automod/models"
)

type FilterInterface interface {
	Name() string
	Args() int
	NewItem(*models.Env, []string) (FilterItemInterface, error)
	Description() string
	Deprecated() bool
}

type FilterItemInterface interface {
	Match(env *models.Env) bool
}
