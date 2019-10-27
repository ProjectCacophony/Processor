package models

import (
	"github.com/jinzhu/gorm"
	"gitlab.com/Cacophony/go-kit/state"
)

type Happening interface {
	Chance(ctx Context) float64 // 0.0-1.0
	Init(repo Repo, state *state.State)
	Do(ctx Context)
}

type Repo interface {
	StoreInteraction(interaction *Interaction) error
}

type Interaction struct {
	gorm.Model
	UserID  string
	Message string
}

func (*Interaction) TableName() string {
	return "rpg_interactions"
}
