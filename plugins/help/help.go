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
			summeryText += " | *(Patrons Only)*"
		}

		summeryText += fmt.Sprintf("\n```%s```", event.Translate(pluginHelp.Description))

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
