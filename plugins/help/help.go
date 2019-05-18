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

	if pluginHelp.Hide {
		event.Respond("help.no-plugin-doc")
		return
	}

	output := fmt.Sprintf("__**%s**__", strings.Title(pluginHelp.Name))

	if len(pluginHelp.PermissionsRequired) > 0 {
		output += fmt.Sprintf(" | Requires **%s**", pluginHelp.PermissionsRequired)
	}

	if containsPatronPermission(pluginHelp.PermissionsRequired) {
		output += " | Patrons Only"
	}

	output += fmt.Sprintf("\n%s", event.Translate(pluginHelp.Description))

	if len(pluginHelp.Commands) == 0 {
		event.Respond(output)
		return
	}

	output += "\n\n"
	commandsList := make([]string, len(pluginHelp.Commands))

	for i, command := range pluginHelp.Commands {

		// commands that are only for bot admins should only output if a bot admin is running the help command
		if containsBotAdminPermission(command.PermissionsRequired) && !event.Has(permissions.BotAdmin) {
			continue
		}

		var commandSummary string
		if command.Name != "" {
			commandSummary += "**" + command.Name + "**\n"
		}

		commandText := event.Prefix() + pluginHelp.Name
		for _, param := range command.Params {
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

			commandText += " " + name
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

		commandsList[i] = commandSummary + "\n"
	}

	output += "__**Commands**__\n"
	output += strings.Join(commandsList, "")

	if displayInChannel {

		event.Respond(output)
	} else {
		if !event.DM() {
			_, err := event.Respond("help.message-sent-to-dm")
			event.Except(err)
		}

		output += fmt.Sprintf("\n\nUse `%shelp %s public` to display the commands in a channel.", event.Prefix(), pluginHelp.Name)
		_, err := event.RespondDM(output)
		event.Except(err)
	}
}
