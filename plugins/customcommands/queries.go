package customcommands

import (
	"gitlab.com/Cacophony/go-kit/events"
	"go.uber.org/zap"
)

func (p *Plugin) getCommandEntries(event *events.Event, commandName string) (entries []Entry) {

	// query entries
	err := p.db.Find(&entries, "name = ? and ((is_user_command = true and user_id = ?) or (is_user_command = false and guild_id = ?))",
		commandName,
		event.UserID,
		event.GuildID,
	).Error
	if err != nil {
		event.Logger().Error("error querying custom commands", zap.Error(err))
	}

	return
}

// func (p *Plugin) getAllAvailableEntries(event *events.Event) (entries []Entry) {

// 	// query entries
// 	err := p.db.Find(&entries, "(is_user_command = true and user_id = ?) or (is_user_command = false and guild_id = ?)",
// 		event.UserID,
// 		event.GuildID,
// 	).Error
// 	if err != nil {
// 		event.Logger().Error("error querying custom commands", zap.Error(err))
// 	}

// 	return
// }

func (p *Plugin) getAllUserEntries(event *events.Event) (entries []Entry) {

	// query entries
	err := p.db.Find(&entries, "is_user_command = true and user_id = ?",
		event.UserID,
	).Error
	if err != nil {
		event.Logger().Error("error querying custom commands", zap.Error(err))
	}

	return
}
func (p *Plugin) getAllServerEntries(event *events.Event) (entries []Entry) {

	// query entries
	err := p.db.Find(&entries, "is_user_command = false and guild_id = ?",
		event.GuildID,
	).Error
	if err != nil {
		event.Logger().Error("error querying custom commands", zap.Error(err))
	}

	return
}

func (p *Plugin) getCommandsByTriggerCount(event *events.Event) (entries []Entry) {

	err := p.db.
		Table((&Entry{}).TableName()).
		Select("name, sum(triggered) as triggered, is_user_command").
		Where("((is_user_command = true and user_id = ?) or (is_user_command = false and guild_id = ?))", event.UserID, event.GuildID).
		Group("name, is_user_command").
		Find(&entries).Error

	if err != nil {
		event.Logger().Error("error querying custom commands", zap.Error(err))
	}
	return
}
