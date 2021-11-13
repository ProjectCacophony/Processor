package customcommands

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	humanize "github.com/dustin/go-humanize"
	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) runCustomCommand(event *events.Event) bool {
	if len(event.Fields()) == 0 {
		return false
	}

	// get all commands with that name
	commands := p.getCommandEntries(event, event.OriginalCommand())
	if len(commands) == 0 {
		return false
	}

	// are we getting the info for this command
	if len(event.Fields()) == 2 && event.Fields()[1] == "info" {
		p.displayCommandsInfo(event, commands)
		return true
	}

	// check if one command was returned, if so run that.
	if len(commands) == 1 {
		if commands[0].IsUserCommand {
			if !p.canUseUserCommand(event) {
				event.Respond("customcommands.cant-use", "level", true)
				return true
			}
		} else {
			if !p.canUseServerCommand(event) {
				event.Respond("customcommands.cant-use", "level", false)
				return true
			}
		}
		err := commands[0].run(event)
		event.Except(err)
		return true
	}

	// multiple commands exist with the given name, get a random one
	userCommands, serverCommands := seporateUserAndServerCommands(commands)

	var command CustomCommand
	var denidedByServerCommandPerm bool
	if len(serverCommands) > 0 {
		if p.canUseServerCommand(event) {
			seed := rand.Intn(len(serverCommands))
			command = serverCommands[seed]
		} else {
			denidedByServerCommandPerm = true
		}
	}

	if len(userCommands) > 0 && command.Name == "" {
		if p.canUseUserCommand(event) {
			seed := rand.Intn(len(userCommands))
			command = userCommands[seed]
		}
	}

	if command.Name != "" {
		err := command.run(event)
		event.Except(err)
		return true
	}

	if denidedByServerCommandPerm {
		event.Respond("customcommands.cant-use", "level", false)
		return true
	}

	event.Respond("customcommands.cant-use", "level", true)
	return true
}

func (p *Plugin) runRandomCommand(event *events.Event) {
	if len(event.Fields()) == 0 {
		return
	}

	var denidedByServerCommandPerm bool
	var denidedByUserCommandPerm bool
	var commands []CustomCommand
	if isUserOperation(event) {
		if p.canUseUserCommand(event) {
			commands = p.getAllUserCommands(event, customCommandTypeContent)
		} else {
			denidedByServerCommandPerm = true
		}
	} else {
		if p.canUseServerCommand(event) {
			commands = p.getAllServerCommands(event, customCommandTypeContent)
		} else {
			denidedByUserCommandPerm = true
		}
	}

	if denidedByServerCommandPerm {
		event.Respond("customcommands.cant-use", "level", false)
		return
	} else if denidedByUserCommandPerm {
		event.Respond("customcommands.cant-use", "level", true)
		return
	}

	if len(commands) == 0 {
		event.Respond("customcommands.no-commands")
		return
	}

	seed := rand.Intn(len(commands))
	commands[seed].run(event)
}

func (p *Plugin) listCommands(event *events.Event) {
	if len(event.Fields()) > 3 {
		event.Respond("common.invalid-params")
		return
	}

	var displayInPublic bool
	if len(event.Fields()) == 3 && event.Fields()[2] == "public" {
		displayInPublic = true
	}

	// get all commands the user has access to, user and on the server
	commands := p.getCommandsByTriggerCount(event)

	listText := createListCommandOutput(event, commands)

	if displayInPublic {
		_, err := event.Respond(listText)
		event.Except(err)
	} else {

		if !event.DM() {
			_, err := event.Respond("common.message-sent-to-dm")
			event.Except(err)
		}

		_, err := event.RespondDM(listText)
		event.Except(err)
	}
}

func (p *Plugin) getCommandInfo(event *events.Event) {
	if len(event.Fields()) != 3 {
		event.Respond("common.invalid-params")
		return
	}

	commands := p.getCommandEntries(event, event.Fields()[2])
	if len(commands) == 0 {
		event.Respond("customcommands.name-not-found")
		return
	}

	p.displayCommandsInfo(event, commands)
}

func (p *Plugin) displayCommandsInfo(event *events.Event, commands []CustomCommand) {
	sort.Slice(commands, func(i, j int) bool {
		return commands[i].Model.CreatedAt.Before(commands[j].Model.CreatedAt)
	})

	totalTriggered := commands[0].Triggered
	command := commands[0]
	content := commands[0].getContent()
	if command.Type == customCommandTypeCommand {
		content = "command alias: " + event.Prefix() + content
	}
	userInfo := "*Unknown*"

	// if multi commands of same name
	if len(commands) > 1 {

		commandCreator := commands[0].UserID
		user, err := event.State().User(commandCreator)
		if err == nil && user != nil {
			userInfo = fmt.Sprintf("%s#%s", user.Username, user.Discriminator)
		}

		totalTriggered = 0
		commandArray := make([]string, len(commands))
		var commandContent string
		for i, cmd := range commands {

			// combind all triggers from each command of same name
			totalTriggered += cmd.Triggered

			// check if all commands of same name were uploaded by same user
			if cmd.UserID != commandCreator {
				userInfo = "*Multiple Users*"
			}

			commandContent = cmd.getContent()
			if cmd.Type == customCommandTypeCommand {
				commandContent = "command alias: " + event.Prefix() + commandContent
			}
			commandArray[i] = fmt.Sprintf("%d) %s", i+1, commandContent)
		}

		content = "__Multiple Commands__\n"
		content += strings.Join(commandArray, "\n")
	} else {

		// get user info
		user, err := event.State().User(command.UserID)
		if err == nil && user != nil {
			userInfo = fmt.Sprintf("%s#%s", user.Username, user.Discriminator)
		}
	}

	embed := &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Title:       fmt.Sprintf("Custom Command: `%s`", command.Name),
			Description: content,
			Fields: []*discordgo.MessageEmbedField{
				{Name: "Author", Value: userInfo},
				{Name: "Times Triggered", Value: humanize.Comma(int64(totalTriggered))},
				{Name: "Created At", Value: fmt.Sprintf("%s UTC", command.Model.CreatedAt.Format(time.ANSIC))},
			},
		},
	}

	event.RespondComplex(embed)
}

func (p *Plugin) searchCommand(event *events.Event) {
	if len(event.Fields()) != 3 {
		event.Respond("common.invalid-params")
		return
	}

	commands := p.searchForCommand(event, event.Fields()[2])
	if len(commands) == 0 {
		event.Respond("customcommands.name-not-found")
		return
	}

	listText := createListCommandOutput(event, commands)

	_, err := event.Respond(listText)
	event.Except(err)
}

func createListCommandOutput(event *events.Event, commands []CustomCommand) (listText string) {
	sort.Slice(commands, func(i, j int) bool {
		return commands[i].Triggered > commands[j].Triggered
	})

	userCommands, serverCommands := seporateUserAndServerCommands(commands)
	var aliasText string

	// List server commands if in a guild
	if event.GuildID != "" {
		guild, err := event.State().Guild(event.GuildID)
		if err != nil {
			event.Except(err)
			return
		}

		serverList := make([]string, len(serverCommands))
		for i, command := range serverCommands {
			aliasText = ""
			if command.Type == customCommandTypeCommand {
				aliasText = "command alias, "
			}
			serverList[i] = fmt.Sprintf("`%s` (%sused %d times)", command.Name, aliasText, command.Triggered)
		}

		if len(serverList) != 0 {
			listText += fmt.Sprintf("Custom Commands on `%s` (%d):\n", guild.Name, len(serverCommands))
			listText += strings.Join(serverList, "\n") + "\n"
		}
	}

	// List user commands
	userList := make([]string, len(userCommands))
	for i, command := range userCommands {
		aliasText = ""
		if command.Type == customCommandTypeCommand {
			aliasText = "command alias, "
		}
		userList[i] = fmt.Sprintf("`%s` (%sused %d times)", command.Name, aliasText, command.Triggered)
	}
	if len(userList) != 0 {
		listText += fmt.Sprintf("\nYour Custom Commands (%d):\n", len(userCommands))
		listText += strings.Join(userList, "\n")
	}

	if listText == "" {
		listText = "customcommands.no-commands-found"
	}

	return
}
