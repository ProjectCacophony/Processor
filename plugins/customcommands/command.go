package customcommands

import (
	"math/rand"

	"gitlab.com/Cacophony/go-kit/events"
	"go.uber.org/zap"
)

func (p *Plugin) ranCustomCommand(event *events.Event) bool {
	if len(event.Fields()) == 0 {
		return false
	}

	// query entries
	var entries []Entry
	err := p.db.Find(&entries, "name = ? and ((is_user_command = true and user_id = ?) or (is_user_command = false and guild_id = ?))",
		event.Fields()[0],
		event.UserID,
		event.GuildID,
	).Error
	if err != nil {
		event.Logger().Error("error querying custom commands", zap.Error(err))
		return false
	}

	if len(entries) == 0 {
		return false
	}

	if len(entries) == 1 {
		_, err = event.Respond(entries[0].Content)
		event.Except(err)
		return true
	}

	var serverEntries []Entry
	var userEntries []Entry
	for _, entry := range entries {
		if entry.IsUserCommand {
			userEntries = append(userEntries, entry)
		}
		serverEntries = append(serverEntries, entry)
	}

	if len(serverEntries) > 0 {
		seed := rand.Intn(len(serverEntries))
		event.Respond(serverEntries[seed].Content)
		return true
	}

	if len(userEntries) > 0 {
		seed := rand.Intn(len(userEntries))
		event.Respond(userEntries[seed].Content)
		return true
	}

	return false
}
