package help

import (
	"fmt"
	"sort"
	"strings"

	"gitlab.com/Cacophony/go-kit/discord"

	"gitlab.com/Cacophony/go-kit/events"
)

const PluginHelpListKey = "pluginHelpList"

func listCommands(event *events.Event, displayInChannel bool) {

	var pluginHelpList []PluginHelp
	var ok bool
	if pluginHelpList, ok = event.Context().Value(PluginHelpListKey).([]PluginHelp); !ok {
		// TODO (snake): this shouldn't happen, need to report some error to sentry
		return
	}

	// sort plugins by name
	sort.Slice(pluginHelpList, func(i, j int) bool {
		return pluginHelpList[i].PluginName < pluginHelpList[j].PluginName
	})

	// build each plugins help text with the plugin name and description
	pluginNames := make([]string, len(pluginHelpList))
	for _, pluginHelp := range pluginHelpList {
		if pluginHelp.Hide {
			continue
		}

		summeryText := fmt.Sprintf("__**%s**__ | `%shelp %s`", strings.Title(pluginHelp.PluginName), event.Prefix(), pluginHelp.PluginName)

		if pluginHelp.PatreonOnly {
			summeryText += " | *(Patrons Only)*"
		}

		summeryText += fmt.Sprintf("\n```%s```", pluginHelp.Description)

		pluginNames = append(pluginNames, summeryText)
	}

	helpText := strings.Join(pluginNames, "\n")

	if displayInChannel {
		event.Respond(helpText)
	} else {
		event.Respond("help.message-sent-to-dm")
		dmChannel, err := discord.DMChannel(event.Redis(), event.Discord(), event.UserID)
		event.Logger().Info(fmt.Sprint("dm channel", dmChannel))
		if err != nil {
			event.Except(err)
			return
		}
		helpText += fmt.Sprintf("\n\nUse `%shelp public` to display the commands in your server", event.Prefix())
		event.Send(dmChannel, helpText, false)
	}

}
