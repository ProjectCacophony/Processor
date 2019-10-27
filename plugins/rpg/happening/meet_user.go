package happening

import (
	"fmt"
	"math/rand"

	"gitlab.com/Cacophony/Processor/plugins/rpg/models"
	"gitlab.com/Cacophony/go-kit/state"
	"go.uber.org/zap"
)

type MeetUser struct {
	repo  models.Repo
	state *state.State
}

func (h *MeetUser) Init(repo models.Repo, state *state.State) {
	h.repo = repo
	h.state = state
}

func (h *MeetUser) Chance(ctx models.Context) float64 {
	for _, mention := range ctx.Message.Mentions {
		if mention.Bot {
			continue
		}

		return 1
	}

	return 0.25
}

func (h *MeetUser) Do(ctx models.Context) {
	var mentionID string
	if len(ctx.Message.Mentions) > 0 {
		mentionID = ctx.Message.Mentions[0].ID
	}

	if mentionID == "" {
		guildMembers, err := h.state.GuildMembers(ctx.GuildID)
		if err != nil {
			zap.L().Error("happening error", zap.Error(err))
			return
		}

		mentionID = guildMembers[rand.Intn(len(guildMembers))]
	}

	if mentionID == "" {
		return
	}

	user, err := h.state.User(mentionID)
	if err != nil {
		zap.L().Error("happening error", zap.Error(err))
		return
	}

	err = h.repo.StoreInteraction(&models.Interaction{
		UserID:  ctx.UserID,
		Message: fmt.Sprintf("You met %s", user.String()),
	})
	if err != nil {
		zap.L().Error("happening error", zap.Error(err))
	}
}
