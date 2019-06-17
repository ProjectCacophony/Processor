package automod

import (
	"gitlab.com/Cacophony/Processor/plugins/automod/models"
	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) cmdHistory(event *events.Event) {
	var logs []models.LogEntry
	err := p.db.
		Preload("Rule").
		Preload("Rule.Filters").
		Preload("Rule.Actions").
		Where("guild_id = ?", event.MessageCreate.GuildID).
		Limit(10).
		Order("created_at DESC").
		Find(&logs).Error
	if err != nil {
		event.Except(err)
		return
	}

	// reverse slice (from oldest to newest)
	for i := len(logs)/2 - 1; i >= 0; i-- {
		opp := len(logs) - 1 - i
		logs[i], logs[opp] = logs[opp], logs[i]
	}

	_, err = event.Respond(
		"automod.history.content",
		"logs", logs,
	)
	event.Except(err)
}
