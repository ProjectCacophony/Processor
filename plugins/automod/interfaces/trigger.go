package interfaces

import (
	"gitlab.com/Cacophony/Processor/plugins/automod/models"
)

type TriggerInterface interface {
	Name() string
	NewItem() TriggerItemInterface
}

type TriggerItemInterface interface {
	Match(env *models.Env) bool
}
