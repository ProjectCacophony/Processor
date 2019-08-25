package help

import (
	"fmt"
	"sort"
	"strings"

	"gitlab.com/Cacophony/Processor/plugins/common"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/permissions"
)

const PluginHelpListKey = "pluginHelpList"

func listCommands(event *events.Event, pluginHelpList []*common.PluginHelp, displayInChannel bool) {

	// sort plugins by name
	sort.Slice(pluginHelpList, func(i, j int) bool {
		return pluginHelpList[i].Name < pluginHelpList[j].Name
	})

	// build each plugins help text with the plugin name and description
	pluginNames := make([]string, len(pluginHelpList))
	for _, pluginHelp := range pluginHelpList {
		if pluginHelp.Hide {
			continue
		}

		summeryText := fmt.Sprintf("__**%s**__ | `%shelp %s`",
			strings.Title(pluginHelp.Name), event.Prefix(), pluginHelp.Name)

		if containsBotAdminPermission(pluginHelp.PermissionsRequired) && !event.Has(permissions.BotAdmin) {
			continue
		}

		if containsPatronPermission(pluginHelp.PermissionsRequired) {
			summeryText += " | Patrons Only"
		}

		summeryText += fmt.Sprintf("\n%s\n", event.Translate(pluginHelp.Description))

		pluginNames = append(pluginNames, summeryText)
	}

	helpText := strings.Join(pluginNames, "\n")

	if displayInChannel {
		_, err := event.Respond(helpText)
		event.Except(err)
	} else {

		if !event.DM() {
			_, err := event.Respond("help.message-sent-to-dm")
			event.Except(err)
		}

		helpText += fmt.Sprintf("\n\nUse `%shelp public` to display the commands in a channel.", event.Prefix())
		_, err := event.RespondDM(helpText)
		event.Except(err)
	}

}

func displayPluginCommands(event *events.Event, pluginHelp *common.PluginHelp, displayInChannel bool) {
	// hidden module commands should only show for bot admins in dm
	if pluginHelp.Hide {

		if !event.Has(permissions.BotAdmin) {
			event.Respond("help.no-plugin-doc")
			return
		}
	}

	// Display module name
	output := fmt.Sprintf("__**%s**__", strings.Title(pluginHelp.Name))

	// Check and display overall module permissions
	if len(pluginHelp.PermissionsRequired) > 0 {
		output += fmt.Sprintf(" | Requires **%s**", pluginHelp.PermissionsRequired)
	}
	if containsPatronPermission(pluginHelp.PermissionsRequired) {
		output += " | Patrons Only"
	}

	// Output module description
	output += fmt.Sprintf("\n%s\n\n", event.Translate(pluginHelp.Description))

	// Get module reactions if any
	reactionList := make([]string, len(pluginHelp.Reactions))
	for i, reaction := range pluginHelp.Reactions {
		if containsBotAdminPermission(reaction.PermissionsRequired) && !event.Has(permissions.BotAdmin) {
			continue
		}
		reactionList[i] = formatReaction(event, reaction)
	}

	// Get module commands if any
	commandsList := make([]string, len(pluginHelp.Commands))
	for i, command := range pluginHelp.Commands {
		if containsBotAdminPermission(command.PermissionsRequired) && !event.Has(permissions.BotAdmin) {
			continue
		}
		commandsList[i] = formatCommand(event, command, pluginHelp.Name)
	}

	if len(reactionList) > 0 {
		output += "__**Reactions**__\n"
		output += strings.Join(reactionList, "") + "\n"
	}

	if len(commandsList) > 0 {
		output += "__**Commands**__\n"
		output += strings.Join(commandsList, "")
	}

	// output to dm or channel
	if displayInChannel {
		event.Respond(output)
	} else {
		if !event.DM() {
			_, err := event.Respond("help.message-sent-to-dm")
			event.Except(err)
		}

		output += fmt.Sprintf("\n\nUse `%s%s help public` to display the commands in a channel.", event.Prefix(), pluginHelp.Name)
		_, err := event.RespondDM(output)
		event.Except(err)
	}
}

func formatReaction(event *events.Event, reaction common.Reaction) string {
	var reactionSummary string
	reactionSummary += fmt.Sprintf(
		"%s **%s**",
		event.Translate(reaction.EmojiName),
		event.Translate(reaction.Description),
	)

	var requirements []string
	if containsPatronPermission(reaction.PermissionsRequired) {
		requirements = append(requirements, "Patrons Only")
	}

	if len(reaction.PermissionsRequired) > 0 {
		requirements = append(requirements, fmt.Sprintf("Requires *%s*", reaction.PermissionsRequired))
	}

	if len(requirements) > 0 {
		reactionSummary += "\n\t\t- " + strings.Join(requirements, " | ")
	}

	return reactionSummary + "\n"
}

func formatCommand(event *events.Event, command common.Command, pluginName string) string {
	var commandSummary string
	if command.Name != "" {
		commandSummary += "**" + event.Translate(command.Name) + "**\n"
	}

	var commandText string
	if !command.SkipPrefix {
		commandText += event.Prefix()
	}
	if !command.SkipRootCommand {
		commandText += pluginName
	}

	for i, param := range command.Params {
		name := param.Name

		if param.Optional {
			name = "?" + name
		}

		switch param.Type {
		case common.Flag:
		case common.QuotedText:
			name = "\"<" + name + ">\""
		case common.Text:
			name = "<" + name + ">"
		case common.User:
			name = "<@" + name + ">"
		case common.Channel:
			name = "<#" + name + ">"

		}

		if command.SkipRootCommand && i == 0 {
			commandText += name
		} else {
			commandText += " " + name
		}
	}
	commandText = "`" + commandText + "`"

	commandSummary += commandText

	if command.Description != "" {
		commandSummary += fmt.Sprintf("\n\t\t*%s*", event.Translate(command.Description))
	}

	var requirements []string

	if containsPatronPermission(command.PermissionsRequired) {
		requirements = append(requirements, "Patrons Only")
	}

	if len(command.PermissionsRequired) > 0 {
		requirements = append(requirements, fmt.Sprintf("Requires *%s*", command.PermissionsRequired))
	}

	if len(requirements) > 0 {
		commandSummary += "\n\t\t- " + strings.Join(requirements, " | ")
	}

	return commandSummary + "\n"
}
