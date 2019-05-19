package customcommands

import (
	"math/rand"

	"gitlab.com/Cacophony/go-kit/events"
	"go.uber.org/zap"
)

func (p *Plugin) runCustomCommand(event *events.Event) bool {
	if len(event.Fields()) == 0 {
		return false
	}

	entries := p.getCommandEntries(event, event.Fields()[0])

	if len(entries) == 0 {
		return false
	}

	if len(entries) == 1 {
		_, err := event.Respond(entries[0].Content)
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

func (p *Plugin) runRandomCommand(event *events.Event) {
	if len(event.Fields()) == 0 {
		return
	}

	var entries []Entry
	if isUserOperation(event) {
		entries = p.getAllUserEntries(event)
	} else {
		entries = p.getAllServerEntries(event)
	}

	if len(entries) == 0 {
		event.Respond("customcommands.no-commands")
		return
	}

	seed := rand.Intn(len(entries))
	event.Respond(entries[seed].Content)
}

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
