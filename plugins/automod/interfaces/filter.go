package interfaces

import (
	"gitlab.com/Cacophony/Processor/plugins/automod/models"
)

type FilterInterface interface {
	Name() string
	NewItem(string) (FilterItemInterface, error)
}

type FilterItemInterface interface {
	Match(env *models.Env) bool
}
