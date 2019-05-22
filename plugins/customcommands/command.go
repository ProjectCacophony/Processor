package customcommands

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"

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

	userEntries, serverEntries := seporateUserAndServerEntries(entries)

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

func (p *Plugin) listCommands(event *events.Event) {
	if len(event.Fields()) > 3 {
		event.Respond("common.invalid-params")
		return
	}

	var displayInPublic bool
	if len(event.Fields()) == 3 && event.Fields()[2] == "public" {
		displayInPublic = true
	}

	var entries []Entry
	err := p.db.Select("name, sum(triggered) as triggered, is_user_command").Table((&Entry{}).TableName()).
		Where("((is_user_command = true and user_id = ?) or (is_user_command = false and guild_id = ?))", event.UserID, event.GuildID).
		Group("name, is_user_command").Find(&entries).Error
	if err != nil {
		event.Respond("customcommands.no-commands")
		return
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name < entries[j].Name
	})

	userEntries, serverEntries := seporateUserAndServerEntries(entries)

	var userList []string
	for _, entry := range userEntries {
		userList = append(userList, fmt.Sprintf("`%s` (used %d times)", entry.Name, entry.Triggered))
	}

	var serverList []string
	for _, entry := range serverEntries {
		serverList = append(serverList, fmt.Sprintf("`%s` (used %d times)", entry.Name, entry.Triggered))
	}

	guild, err := event.State().Guild(event.GuildID)
	event.Except(err)

	listText := fmt.Sprintf("Custom Commands on `%s` (%d):\n", guild.Name, len(serverEntries))
	listText += strings.Join(serverList, "\n")

	listText += fmt.Sprintf("\n\nYour Custom Commands (%d):\n", len(userEntries))
	listText += strings.Join(userList, "\n")

	if displayInPublic {
		_, err = event.Respond(listText)
		event.Except(err)
	} else {

		if !event.DM() {
			_, err := event.Respond("common.message-sent-to-dm")
			event.Except(err)
		}

		_, err = event.RespondDM(listText)
		event.Except(err)
	}
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

func (p *Plugin) getAllAvailableEntries(event *events.Event) (entries []Entry) {

	// query entries
	err := p.db.Find(&entries, "(is_user_command = true and user_id = ?) or (is_user_command = false and guild_id = ?)",
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

func seporateUserAndServerEntries(entries []Entry) (userEntries []Entry, serverEntries []Entry) {
	for _, entry := range entries {
		if entry.IsUserCommand {
			userEntries = append(userEntries, entry)
		} else {
			serverEntries = append(serverEntries, entry)
		}
	}
	return
}
