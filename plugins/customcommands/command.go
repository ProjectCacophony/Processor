package customcommands

import (
	"gitlab.com/Cacophony/go-kit/events"
	"go.uber.org/zap"
)

func (p *Plugin) ranCustomCommand(event *events.Event) bool {
	if len(event.Fields()) == 0 {
		return false
	}

	event.Logger().Info(event.Fields()[0])

	var entires []Entry
	err := p.db.Find(&entires, "name = ? and ((is_user_command = true and user_id = ?) or (is_user_command = false and guild_id = ?))",
		event.Fields()[0],
		event.UserID,
		event.GuildID,
	).Order("is_user_command").Error
	if err != nil {
		event.Logger().Error("error querying custom commands", zap.Error(err))
		return false
	}

	if len(entires) == 0 {
		return false
	}

	event.Respond(entires[0].Content)

	return true
}
