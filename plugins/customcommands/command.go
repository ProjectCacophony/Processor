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
		if entries[0].IsUserCommand {
			if !p.canUseUserCommand(event) {
				event.Respond("customcommands.cant-use", "level", true)
				return true
			}
		} else {
			if !p.canUseServerCommand(event) {
				event.Respond("customcommands.cant-use", "level", false)
				return true
			}
		}
		err := entries[0].run(event)
		event.Except(err)
		return true
	}

	userEntries, serverEntries := seporateUserAndServerEntries(entries)

	var entry Entry
	var denidedByServerCommandPerm bool
	var denidedByUserCommandPerm bool
	if len(serverEntries) > 0 {
		if p.canUseServerCommand(event) {
			seed := rand.Intn(len(serverEntries))
			entry = serverEntries[seed]
		} else {
			denidedByServerCommandPerm = true
		}
	}

	if len(userEntries) > 0 && entry == (Entry{}) {
		if p.canUseUserCommand(event) {
			seed := rand.Intn(len(userEntries))
			entry = userEntries[seed]
		} else {
			denidedByUserCommandPerm = true
		}
	}

	if denidedByServerCommandPerm {
		event.Respond("customcommands.cant-use", "level", false)
		return true
	} else if denidedByUserCommandPerm {
		event.Respond("customcommands.cant-use", "level", true)
		return true
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

	var denidedByServerCommandPerm bool
	var denidedByUserCommandPerm bool
	var entries []Entry
	if isUserOperation(event) {
		if p.canUseUserCommand(event) {
			entries = p.getAllUserEntries(event)
		} else {
			denidedByServerCommandPerm = true
		}
	} else {
		if p.canUseServerCommand(event) {
			entries = p.getAllServerEntries(event)
		} else {
			denidedByUserCommandPerm = true
		}
	}

	if denidedByServerCommandPerm {
		event.Respond("customcommands.cant-use", "level", false)
		return
	} else if denidedByUserCommandPerm {
		event.Respond("customcommands.cant-use", "level", true)
		return
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

		serverList := make([]string, len(serverEntries))
		for i, entry := range serverEntries {
			serverList[i] = fmt.Sprintf("`%s` (used %d times)", entry.Name, entry.Triggered)
		}

		listText += fmt.Sprintf("Custom Commands on `%s` (%d):\n", guild.Name, len(serverEntries))
		listText += strings.Join(serverList, "\n")
	}

	// List user commands
	userList := make([]string, len(userEntries))
	for i, entry := range userEntries {
		userList[i] = fmt.Sprintf("`%s` (used %d times)", entry.Name, entry.Triggered)
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
