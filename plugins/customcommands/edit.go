package customcommands

import (
	"fmt"

	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) addCommand(event *events.Event) {
	if len(event.Fields()) < 4 {
		event.Respond("common.invalid-params")
		return
	}

	// TODO: check for attachments

	if !p.canEditCommand(event) {
		event.Respond("customcommands.cant-edit")
		return
	}

	var cmdName string
	var cmdContent string
	if hasUserParam(event) {
		cmdName = event.Fields()[3]
		cmdContent = event.Fields()[4]
	} else {
		cmdName = event.Fields()[2]
		cmdContent = event.Fields()[3]
	}

	if !isValidCommandName(cmdName) {
		event.Respond("customcommands.name.no-spaces")
		return
	}

	err := entryAdd(
		p.db,
		cmdName,
		event.UserID,
		event.GuildID,
		cmdContent,
		isUserOperation(event),
	)

	if err != nil {
		event.Except(err)
		return
	}

	_, err = event.Respond("customcommands.add.success",
		"cmdName", cmdName,
		"isUserCommand", isUserOperation(event),
	)
	event.Except(err)
}

func (p *Plugin) editCommand(event *events.Event) {
	if len(event.Fields()) < 4 {
		event.Respond("common.invalid-params")
		return
	}

	if !p.canEditCommand(event) {
		event.Respond("customcommands.cant-edit")
		return
	}

	var cmdName string
	var cmdContent string
	if hasUserParam(event) {
		cmdName = event.Fields()[3]
		cmdContent = event.Fields()[4]
	} else {
		cmdName = event.Fields()[2]
		cmdContent = event.Fields()[3]
	}

	if !isValidCommandName(cmdName) {
		event.Respond("customcommands.name.no-spaces")
		return
	}

	tempEntries := p.getCommandEntries(event, cmdName)

	// filter out server or user entries
	var entries []Entry
	isUserOp := isUserOperation(event)
	for _, entry := range tempEntries {
		if entry.IsUserCommand == isUserOp {
			entries = append(entries, entry)
		}
	}

	var err error
	switch len(entries) {
	case 0:
		event.Respond("customcommands.name-not-found")
		return
	case 1:
		entries[0].Content = cmdContent
		err = entryUpsert(p.db, &entries[0])
	default:
		// TODO: finish processing the one they want to update
		askEntryToUpdate(event, entries)
		return
	}

	if err != nil {
		event.Except(err)
		return
	}

	_, err = event.Respond("customcommands.update.success",
		"cmdName", cmdName,
		"isUserCommand", isUserOperation(event),
	)
	event.Except(err)
}

func askEntryToUpdate(event *events.Event, entries []Entry) {
	output := "**Which command would you like to update?**```"

	for i, entry := range entries {
		output += fmt.Sprintf("%d: %s\n", i, entry.Content)
	}
	output += "```"

	event.Respond(output)
}

func (p *Plugin) deleteCommand(event *events.Event) {
	if len(event.Fields()) < 3 {
		event.Respond("common.invalid-params")
		return
	}

	if !p.canEditCommand(event) {
		event.Respond("customcommands.cant-edit")
		return
	}

	var cmdName string
	if hasUserParam(event) {
		cmdName = event.Fields()[3]
	} else {
		cmdName = event.Fields()[2]
	}

	if !isValidCommandName(cmdName) {
		event.Respond("customcommands.name.no-spaces")
		return
	}

	tempEntries := p.getCommandEntries(event, cmdName)

	// filter out server or user entries
	var entries []Entry
	isUserOp := isUserOperation(event)
	for _, entry := range tempEntries {
		if entry.IsUserCommand == isUserOp {
			entries = append(entries, entry)
		}
	}

	var err error
	switch len(entries) {
	case 0:
		event.Respond("customcommands.name-not-found")
		return
	case 1:
		err = entryRemove(p.db, entries[0].Model.ID)
	default:
		// TODO: finish processing the one they want to delete
		askEntryToDelete(event, entries)
		return
	}

	if err != nil {
		event.Except(err)
		return
	}

	_, err = event.Respond("customcommands.delete.success",
		"cmdName", cmdName,
		"isUserCommand", isUserOperation(event),
	)
	event.Except(err)
}
func askEntryToDelete(event *events.Event, entries []Entry) {
	output := "**Which command would you like to delete?**```"

	for i, entry := range entries {
		output += fmt.Sprintf("%d: %s\n", i, entry.Content)
	}
	output += "```"

	event.Respond(output)
}
