package customcommands

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"

	"gitlab.com/Cacophony/go-kit/events"
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
		err := entries[0].run(event)
		event.Except(err)
		return true
	}

	userEntries, serverEntries := seporateUserAndServerEntries(entries)

	var entry Entry
	if len(serverEntries) > 0 {
		seed := rand.Intn(len(serverEntries))
		entry = serverEntries[seed]
	}

	if len(userEntries) > 0 {
		seed := rand.Intn(len(userEntries))
		entry = userEntries[seed]
	}

	if entry != (Entry{}) {
		err := entry.run(event)
		event.Except(err)
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
	entries[seed].run(event)
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

	// get all commands the user has access to, user and on the server
	entries := p.getCommandsByTriggerCount(event)
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Triggered > entries[j].Triggered
	})

	var listText string
	userEntries, serverEntries := seporateUserAndServerEntries(entries)

	// List server commands if in a guild
	if event.GuildID != "" {
		guild, err := event.State().Guild(event.GuildID)
		if err != nil {
			event.Except(err)
			return
		}

		var serverList []string
		for _, entry := range serverEntries {
			serverList = append(serverList, fmt.Sprintf("`%s` (used %d times)", entry.Name, entry.Triggered))
		}

		listText += fmt.Sprintf("Custom Commands on `%s` (%d):\n", guild.Name, len(serverEntries))
		listText += strings.Join(serverList, "\n")
	}

	// List user commands
	var userList []string
	for _, entry := range userEntries {
		userList = append(userList, fmt.Sprintf("`%s` (used %d times)", entry.Name, entry.Triggered))
	}
	listText += fmt.Sprintf("\n\nYour Custom Commands (%d):\n", len(userEntries))
	listText += strings.Join(userList, "\n")

	if displayInPublic {
		_, err := event.Respond(listText)
		event.Except(err)
	} else {

		if !event.DM() {
			_, err := event.Respond("common.message-sent-to-dm")
			event.Except(err)
		}

		_, err := event.RespondDM(listText)
		event.Except(err)
	}
}
