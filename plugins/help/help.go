package help

import (
	"fmt"
	"sort"
	"strings"

	"gitlab.com/Cacophony/Processor/plugins/common"
	"gitlab.com/Cacophony/go-kit/events"
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

	if len(pluginHelp.ParamSets) == 0 {
		event.Respond(output)
		return
	}

	output += "\n\n"
	commandsList := make([]string, len(pluginHelp.ParamSets))

	for i, paramSet := range pluginHelp.ParamSets {
		command := event.Prefix() + pluginHelp.Name

		for _, param := range paramSet.Params {
			name := param.Name

			if param.Optional {
				name = "?" + name
			}

			if !param.NotVariable {
				name = "`<" + name + ">`"
			}

			command += " " + name
		}

		if paramSet.Description != "" {
			command += fmt.Sprintf("\n\t*%s*", paramSet.Description)
		}

		var requirements []string

		if paramSet.PatreonOnly {
			requirements = append(requirements, "Patrons Only")
		}

		if len(paramSet.PermissionsRequired) > 0 {
			requirements = append(requirements, fmt.Sprintf("Requires **%s**", paramSet.PermissionsRequired))
		}

		if len(requirements) > 0 {
			command += "\n\t- " + strings.Join(requirements, " | ")
		}

		commandsList[i] = command
	}

	output += "__**Commands**__\n"
	output += strings.Join(commandsList, "\n")
	event.Respond(output)
}
