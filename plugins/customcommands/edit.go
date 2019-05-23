package customcommands

import (
	"encoding/json"
	"errors"
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

	switch len(entries) {
	case 0:
		event.Respond("customcommands.name-not-found")
		return
	case 1:
		entries[0].Content = cmdContent
		p.processCommandEdit(event, entries[0], isUserOperation(event))
	default:
		output := "**Multiple commands with this name, which command would you like to update?**```"
		for i, entry := range entries {
			output += fmt.Sprintf("%d: %s\n", i+1, entry.Content)
		}
		output += "```"
		event.Respond(output)

		openEditQuestionnaire(event, entries, cmdContent, isUserOperation(event), true)
		return
	}
}

func (p *Plugin) handleEditResponse(event *events.Event, enteredNum int) bool {

	newContent, ok := event.QuestionnaireMatch.Payload["newContent"].(string)
	if !ok {
		event.Except(errors.New("received invalid questionnaire match payload"))
		return true
	}

	commandBytes, ok := event.QuestionnaireMatch.Payload["commands"].(string)
	if !ok {
		event.Except(errors.New("received invalid questionnaire match payload"))
		return true
	}

	var commands []Entry
	err := json.Unmarshal([]byte(commandBytes), &commands)
	event.Except(err)

	isUserOperation, ok := event.QuestionnaireMatch.Payload["isUserOperation"].(bool)
	if !ok {
		event.Except(errors.New("received invalid questionnaire match payload"))
		return true
	}

	if len(commands) < enteredNum {
		openEditQuestionnaire(event, commands, newContent, isUserOperation, true)
		_, err = event.Send(
			event.ChannelID,
			"Invalid number entered.",
		)
		event.Except(err)
		return false
	}

	command := commands[enteredNum-1]
	command.Content = newContent

	p.processCommandEdit(event, command, isUserOperation)
	return true
}

func (p *Plugin) processCommandEdit(event *events.Event, command Entry, isUserOperation bool) {

	err := entryUpsert(p.db, &command)
	if err != nil {
		event.Except(err)
		return
	}

	_, err = event.Send(
		event.ChannelID,
		"customcommands.update.success",
		"cmdName", command.Name,
		"isUserCommand", isUserOperation,
	)
	event.Except(err)
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

	switch len(entries) {
	case 0:
		event.Respond("customcommands.name-not-found")
		return
	case 1:
		p.processCommandDelete(event, entries[0], isUserOp)
	default:
		output := "**Multiple commands with this name, which command would you like to delete?**```"
		for i, entry := range entries {
			output += fmt.Sprintf("%d: %s\n", i+1, entry.Content)
		}
		output += "```"

		event.Respond(output)
		openEditQuestionnaire(event, entries, "", isUserOperation(event), false)
		return
	}
}

func (p *Plugin) handleDeleteResponse(event *events.Event, enteredNum int) bool {

	commandBytes, ok := event.QuestionnaireMatch.Payload["commands"].(string)
	if !ok {
		event.Except(errors.New("received invalid questionnaire match payload"))
		return true
	}

	var commands []Entry
	err := json.Unmarshal([]byte(commandBytes), &commands)
	event.Except(err)

	isUserOperation, ok := event.QuestionnaireMatch.Payload["isUserOperation"].(bool)
	if !ok {
		event.Except(errors.New("received invalid questionnaire match payload"))
		return true
	}

	if len(commands) < enteredNum {
		openEditQuestionnaire(event, commands, "", isUserOperation, false)
		_, err = event.Send(
			event.ChannelID,
			"Invalid number entered.",
		)
		event.Except(err)
		return true
	}

	p.processCommandDelete(event, commands[enteredNum-1], isUserOperation)
	return true
}

func (p *Plugin) processCommandDelete(event *events.Event, command Entry, isUserOperation bool) {

	err := entryRemove(p.db, command.Model.ID)
	if err != nil {
		event.Except(err)
		return
	}

	_, err = event.Send(
		event.ChannelID,
		"customcommands.delete.success",
		"cmdName", command.Name,
		"isUserCommand", isUserOperation,
	)
	event.Except(err)
}

func openEditQuestionnaire(event *events.Event, commands []Entry, newContent string, isUserOp, isEdit bool) {
	commandBytes, err := json.Marshal(commands)
	event.Except(err)

	payload := map[string]interface{}{
		"commands":        string(commandBytes),
		"isUserOperation": isUserOp,
	}

	key := deleteQuestionnaireKey
	if isEdit {
		key = editQuestionnaireKey
		payload["newContent"] = newContent
	}

	err = event.Questionnaire().Register(
		key,
		events.QuestionnaireFilter{
			GuildID:   event.GuildID,
			ChannelID: event.ChannelID,
			UserID:    event.UserID,
			Type:      events.MessageCreateType,
		},
		payload,
	)
	if err != nil {
		event.Except(err)
		return
	}
}
