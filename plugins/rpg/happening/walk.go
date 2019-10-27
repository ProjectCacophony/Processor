package happening

import (
	"gitlab.com/Cacophony/Processor/plugins/rpg/models"
	"gitlab.com/Cacophony/go-kit/state"
	"go.uber.org/zap"
)

type Walk struct {
	repo  models.Repo
	state *state.State
}

func (h *Walk) Init(repo models.Repo, state *state.State) {
	h.repo = repo
	h.state = state
}

func (h *Walk) Chance(ctx models.Context) float64 {
	return 1
}

func (h *Walk) Do(ctx models.Context) {
	err := h.repo.StoreInteraction(&models.Interaction{
		UserID:  ctx.UserID,
		Message: "You walked around",
	})
	if err != nil {
		zap.L().Error("happening error", zap.Error(err))
	}
}
