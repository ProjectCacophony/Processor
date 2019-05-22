package customcommands

import (
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
