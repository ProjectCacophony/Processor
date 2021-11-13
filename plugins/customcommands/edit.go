package customcommands

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) addCommand(event *events.Event, ccType customCommandType) {
	if len(event.Fields()) < 4 && (len(event.Fields()) < 3 && len(event.MessageCreate.Attachments) == 0) {
		event.Respond("common.invalid-params")
		return
	}

	if !p.canEditCommand(event) {
		event.Respond("customcommands.cant-edit")
		return
	}

	var commandFile *events.FileInfo
	if ccType == customCommandTypeContent && event.MessageCreate != nil && len(event.MessageCreate.Attachments) > 0 {

		msgs, err := event.Respond("common.uploading-file")
		if err != nil {
			event.Except(err)
			return
		}

		commandFile, err = event.AddAttachement(event.MessageCreate.Attachments[0])
		if err != nil {
			event.Discord().Client.ChannelMessageDelete(event.ChannelID, msgs[0].ID)
			event.Except(err)
			return
		}

		event.Discord().Client.ChannelMessageDelete(event.ChannelID, msgs[0].ID)
	}

	var cmdName string
	var cmdContent string
	if hasUserParam(event) {
		cmdName = event.Fields()[3]
		if len(event.Fields()) >= 5 {
			cmdContent = event.Fields()[4]
		}
	} else {
		cmdName = event.Fields()[2]
		if len(event.Fields()) >= 4 {
			cmdContent = event.Fields()[3]
		}
	}

	if !isValidCommandName(cmdName) {
		event.Respond("customcommands.name.no-spaces")
		return
	}

	if ccType == customCommandTypeCommand {
		cmdContent = strings.TrimPrefix(cmdContent, event.Prefix())

		if cmdContent == "" {
			event.Send(event.ChannelID, "customcommands.alias-needs-content")
			return
		}
	}

	err := createCustomCommand(event.DB(), CustomCommand{
		Name:          cmdName,
		UserID:        event.UserID,
		GuildID:       event.GuildID,
		Content:       cmdContent,
		IsUserCommand: isUserOperation(event),
		File:          commandFile,
		Type:          ccType,
	})
	if err != nil {
		if strings.Contains(err.Error(), noContent) {
			event.Respond("customcommands.empty")
			return
		}
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
	if len(event.Fields()) < 4 && (len(event.Fields()) < 3 && len(event.MessageCreate.Attachments) == 0) {
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
		if len(event.Fields()) >= 5 {
			cmdContent = event.Fields()[4]
		}
	} else {
		cmdName = event.Fields()[2]
		if len(event.Fields()) >= 4 {
			cmdContent = event.Fields()[3]
		}
	}

	if !isValidCommandName(cmdName) {
		event.Respond("customcommands.name.no-spaces")
		return
	}

	tempCommands := p.getCommandEntries(event, cmdName)

	// filter out server or user entries
	var commands []CustomCommand
	isUserOp := isUserOperation(event)
	for _, entry := range tempCommands {
		if entry.IsUserCommand == isUserOp {
			commands = append(commands, entry)
		}
	}

	switch len(commands) {
	case 0:
		event.Respond("customcommands.name-not-found")
		return
	case 1:
		var newAttachement *discordgo.MessageAttachment
		if event.MessageCreate != nil && len(event.MessageCreate.Attachments) > 0 {
			newAttachement = event.MessageCreate.Attachments[0]
		}
		p.processCommandEdit(event, commands[0], cmdContent, newAttachement, isUserOperation(event))
	default:
		output := "**Multiple commands with this name, which command would you like to update?**```"

		for i, entry := range commands {
			output += fmt.Sprintf("%d) %s\n", i+1, entry.getContent())
		}

		output += "```"
		event.Respond(output)

		openEditQuestionnaire(event, commands, cmdContent, isUserOperation(event), true)
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

	originalMessageID, ok := event.QuestionnaireMatch.Payload["messageID"].(string)
	if !ok {
		event.Except(errors.New("received invalid questionnaire match payload"))
		return true
	}

	originalMessage, err := event.Discord().Client.ChannelMessage(event.ChannelID, originalMessageID)
	if err != nil {
		return false
	}

	var commands []CustomCommand
	err = json.Unmarshal([]byte(commandBytes), &commands)
	event.Except(err)

	isUserOperation, ok := event.QuestionnaireMatch.Payload["isUserOperation"].(bool)
	if !ok {
		event.Except(errors.New("received invalid questionnaire match payload"))
		return true
	}

	if len(commands) < enteredNum {
		_, err = event.Send(
			event.ChannelID,
			"Invalid number entered.",
		)
		event.Except(err)
		return false
	}

	var newAttachement *discordgo.MessageAttachment
	if originalMessage != nil && len(originalMessage.Attachments) > 0 {
		newAttachement = originalMessage.Attachments[0]
	}

	p.processCommandEdit(event, commands[enteredNum-1], newContent, newAttachement, isUserOperation)
	return true
}

func (p *Plugin) processCommandEdit(event *events.Event, originalCommand CustomCommand, newContent string, newAttachement *discordgo.MessageAttachment, isUserOperation bool) {
	if originalCommand.File != nil {
		err := event.DeleteFile(originalCommand.File)
		if err != nil {
			event.Except(err)
			return
		}
		originalCommand.File = nil
	}

	if originalCommand.Type == customCommandTypeContent && newAttachement != nil {
		msgs, messageErr := event.Send(event.ChannelID, "common.uploading-file")

		newFile, err := event.AddAttachement(newAttachement)
		if err != nil {
			event.Discord().Client.ChannelMessageDelete(event.ChannelID, msgs[0].ID)
			event.Except(err)
			return
		}

		if messageErr == nil && msgs[0] != nil {
			event.Discord().Client.ChannelMessageDelete(event.ChannelID, msgs[0].ID)
		}
		originalCommand.File = newFile
	}

	if originalCommand.Type == customCommandTypeCommand {
		newContent = strings.TrimPrefix(newContent, event.Prefix())

		if newContent == "" {
			event.Send(event.ChannelID, "customcommands.alias-needs-content")
			return
		}
	}

	originalCommand.Content = newContent

	err := upsertCustomCommand(p.db, &originalCommand)
	if err != nil {
		event.Except(err)
		return
	}

	_, err = event.Send(
		event.ChannelID,
		"customcommands.update.success",
		"cmdName", originalCommand.Name,
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

	tempCommands := p.getCommandEntries(event, cmdName)

	// filter out server or user commands
	var commands []CustomCommand
	isUserOp := isUserOperation(event)
	for _, entry := range tempCommands {
		if entry.IsUserCommand == isUserOp {
			commands = append(commands, entry)
		}
	}

	switch len(commands) {
	case 0:
		event.Respond("customcommands.name-not-found")
		return
	case 1:
		p.processCommandDelete(event, commands[0], isUserOp)
	default:
		output := "**Multiple commands with this name, which command would you like to delete?**```"
		for i, entry := range commands {
			output += fmt.Sprintf("%d: %s\n", i+1, entry.getContent())
		}
		output += "```"

		event.Respond(output)
		openEditQuestionnaire(event, commands, "", isUserOperation(event), false)
		return
	}
}

func (p *Plugin) handleDeleteResponse(event *events.Event, enteredNum int) bool {
	commandBytes, ok := event.QuestionnaireMatch.Payload["commands"].(string)
	if !ok {
		event.Except(errors.New("received invalid questionnaire match payload"))
		return true
	}

	var commands []CustomCommand
	err := json.Unmarshal([]byte(commandBytes), &commands)
	event.Except(err)

	isUserOperation, ok := event.QuestionnaireMatch.Payload["isUserOperation"].(bool)
	if !ok {
		event.Except(errors.New("received invalid questionnaire match payload"))
		return true
	}

	if len(commands) < enteredNum {
		_, err = event.Send(
			event.ChannelID,
			"Invalid number entered.",
		)
		event.Except(err)
		return false
	}

	p.processCommandDelete(event, commands[enteredNum-1], isUserOperation)
	return true
}

func (p *Plugin) processCommandDelete(event *events.Event, command CustomCommand, isUserOperation bool) {
	if command.File != nil {
		err := event.DeleteFile(command.File)
		if err != nil {
			event.Except(err)
			return
		}
	}

	err := removeCustomCommand(p.db, &command)
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

func openEditQuestionnaire(event *events.Event, commands []CustomCommand, newContent string, isUserOp, isEdit bool) {
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
		payload["messageID"] = event.MessageID
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
