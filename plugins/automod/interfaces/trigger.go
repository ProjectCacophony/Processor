package interfaces

import (
	"gitlab.com/Cacophony/Processor/plugins/automod/models"
)

type TriggerInterface interface {
	Name() string
	Args() int
	NewItem(*models.Env, []string) (TriggerItemInterface, error)
	Description() string
}

type TriggerItemInterface interface {
	Match(env *models.Env) bool
}
