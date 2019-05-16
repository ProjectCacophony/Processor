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

		if pluginHelp.PatreonOnly {
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
		_, err := event.RespondDM(helpText, false)
		event.Except(err)
	}

}

func displayPluginCommands(event *events.Event, pluginHelp *common.PluginHelp) {

	if pluginHelp.Hide {
		event.Respond("help.no-plugin-doc")
		return
	}

	output := fmt.Sprintf("__**%s**__", strings.Title(pluginHelp.Name))

	if len(pluginHelp.PermissionsRequired) > 0 {
		output += fmt.Sprintf(" | Requires **%s**", pluginHelp.PermissionsRequired)
	}

	if pluginHelp.PatreonOnly {
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
		if strings.Contains(command.PermissionsRequired.String(), "Bot Admin") && !event.Has(permissions.BotAdmin) {
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
			case common.Hardcoded:
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

		if command.PatreonOnly {
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

	output += "__**Commands**__\n\n"
	output += strings.Join(commandsList, "")
	event.Respond(output)
}
