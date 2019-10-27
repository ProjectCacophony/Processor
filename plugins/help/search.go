package help

import (
	"strings"

	"gitlab.com/Cacophony/Processor/plugins/common"
	"gitlab.com/Cacophony/go-kit/events"
	"gitlab.com/Cacophony/go-kit/permissions"
)

const resultLimit = 5

type enhancedCommand struct {
	common.Command
	Text string
}

func (p *Plugin) searchCommands(event *events.Event) {
	searchText := event.FieldsVariadic(2)
	if searchText == "" {
		event.Respond("help.search.too-few")
		return
	}

	commands := make([]enhancedCommand, 0, 10)

	for _, plugin := range p.pluginHelpList {
		if containsBotAdminPermission(plugin.PermissionsRequired) && !event.Has(permissions.BotAdmin) {
			continue
		}

		if plugin.Hide {
			continue
		}

		for _, command := range plugin.Commands {
			command := command

			if containsBotAdminPermission(command.PermissionsRequired) && !event.Has(permissions.BotAdmin) {
				continue
			}

			if !matches(searchText, commandSearchText(event, plugin, &command)) {
				continue
			}

			commands = append(commands, enhancedCommand{
				Command: command,
				Text:    formatCommand(event, command, plugin.Names[0]),
			})

			if len(commands) >= resultLimit {
				break
			}
		}

		if len(commands) >= resultLimit {
			break
		}
	}

	event.Respond("help.search.result", "commands", commands, "search", searchText)
}

func commandSearchText(event *events.Event, plugin *common.PluginHelp, command *common.Command) string {
	return event.Translate(command.Name) + " " +
		event.Translate(command.Description) + " " +
		strings.Join(plugin.Names, " ")
}

func matches(needle string, haystack string) bool {
	haystack = strings.ToLower(haystack)

	for _, needleItem := range strings.Split(strings.ToLower(needle), " ") {
		needleItem = strings.TrimSpace(needleItem)
		if !strings.Contains(haystack, needleItem) {
			return false
		}
	}

	return true
}
