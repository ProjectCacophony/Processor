package automod

import (
	"strings"

	"gitlab.com/Cacophony/go-kit/events"
)

func cmdRemove(event *events.Event) {
	if len(event.Fields()) < 3 {
		event.Respond("automod.remove.too-few") // nolint: errcheck
		return
	}

	var ruleToDelete Rule
	err := event.DB().Where("guild_id = ? AND name = ?",
		event.MessageCreate.GuildID, event.Fields()[2]).First(&ruleToDelete).Error
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			event.Respond("automod.remove.not-found") // nolint: errcheck
			return
		}
		event.Except(err)
		return
	}

	err = event.DB().Delete(&ruleToDelete).Error
	if err != nil {
		event.Except(err)
		return
	}

	_, err = event.Respond("automod.remove.success", "rule", ruleToDelete)
	event.Except(err)
}
