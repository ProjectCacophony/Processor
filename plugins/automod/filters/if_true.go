package filters

import (
	"gitlab.com/Cacophony/Processor/plugins/automod/models"

	"gitlab.com/Cacophony/Processor/plugins/automod/interfaces"
)

type True struct {
}

func (f True) Name() string {
	return "if_true"
}

func (f True) Args() int {
	return 0
}

func (f True) Deprecated() bool {
	return false
}

func (f True) NewItem(env *models.Env, args []string) (interfaces.FilterItemInterface, error) {
	return &TrueItem{}, nil
}

func (f True) Description() string {
	return "automod.filters.if_true"
}

type TrueItem struct {
}

func (f *TrueItem) Match(env *models.Env) bool {
	return env.Event != nil
}
