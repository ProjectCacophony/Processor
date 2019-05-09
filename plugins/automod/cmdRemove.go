package automod

import (
	"strings"

	"gitlab.com/Cacophony/Processor/plugins/automod/models"
	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) cmdRemove(event *events.Event) {
	if len(event.Fields()) < 3 {
		event.Respond("automod.remove.too-few")
		return
	}

	var ruleToDelete models.Rule
	err := p.db.Where("guild_id = ? AND name = ?",
		event.MessageCreate.GuildID, event.Fields()[2]).First(&ruleToDelete).Error
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			event.Respond("automod.remove.not-found")
			return
		}
		event.Except(err)
		return
	}

	err = p.db.Delete(&ruleToDelete).Error
	if err != nil {
		event.Except(err)
		return
	}

	_, err = event.Respond("automod.remove.success", "rule", ruleToDelete)
	event.Except(err)
}
