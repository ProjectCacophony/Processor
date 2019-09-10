package customcommands

import (
	"gitlab.com/Cacophony/go-kit/events"
	"go.uber.org/zap"
)

func (p *Plugin) getCommandEntries(event *events.Event, commandName string) (commands []CustomCommand) {

	// query commands
	err := p.db.
		Preload("File").
		Where("name = ? and ((is_user_command = true and user_id = ?) or (is_user_command = false and guild_id = ?))",
			commandName,
			event.UserID,
			event.GuildID,
		).Find(&commands).Error
	if err != nil {
		event.Logger().Error("error querying custom commands", zap.Error(err))
	}

	return
}

func (p *Plugin) searchForCommand(event *events.Event, searchTerm string) (commands []CustomCommand) {
	// query commands
	err := p.db.Find(&commands, "name like ? and ((is_user_command = true and user_id = ?) or (is_user_command = false and guild_id = ?))",
		"%"+searchTerm+"%",
		event.UserID,
		event.GuildID,
	).Error
	if err != nil {
		event.Logger().Error("error querying custom commands", zap.Error(err))
	}

	return
}

func (p *Plugin) getAllUserCommands(event *events.Event, firstType customCommandType, types ...customCommandType) (commands []CustomCommand) {

	// query commands
	err := p.db.Find(&commands, "is_user_command = true and user_id = ? and type IN (?)",
		event.UserID, append(types, firstType),
	).Error
	if err != nil {
		event.Logger().Error("error querying custom commands", zap.Error(err))
	}

	return
}
func (p *Plugin) getAllServerCommands(event *events.Event, firstType customCommandType, types ...customCommandType) (commands []CustomCommand) {
	// query commands
	err := p.db.Find(&commands, "is_user_command = false and guild_id = ? and type IN (?)",
		event.GuildID, append(types, firstType),
	).Error
	if err != nil {
		event.Logger().Error("error querying custom commands", zap.Error(err))
	}

	return
}

func (p *Plugin) getCommandsByTriggerCount(event *events.Event) (commands []CustomCommand) {

	err := p.db.
		Table((&CustomCommand{}).TableName()).
		Select("name, type, sum(triggered) as triggered, is_user_command").
		Where("((is_user_command = true and user_id = ?) or (is_user_command = false and guild_id = ?))", event.UserID, event.GuildID).
		Group("name, type, is_user_command").
		Find(&commands).Error

	if err != nil {
		event.Logger().Error("error querying custom commands", zap.Error(err))
	}
	return
}
