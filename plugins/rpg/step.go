package rpg

import (
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/Processor/plugins/rpg/models"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (p *Plugin) processStep(guildID, userID string, message *discordgo.Message) {
	randValue := rand.Float64()

	context := models.Context{
		GuildID: guildID,
		UserID:  userID,
		Message: message,
	}

	possibleHappenings := make([]models.Happening, 0)

	for _, happening := range happenings {
		if happening.Chance(context) < randValue {
			continue
		}

		possibleHappenings = append(possibleHappenings, happening)
	}

	if len(possibleHappenings) <= 0 {
		return
	}

	happening := possibleHappenings[rand.Intn(len(possibleHappenings))]

	happening.Do(context)
}
