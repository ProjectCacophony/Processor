package customcommands

const (
	everyoneCreatePermissionKey        = "cacophony:processor:cc:everyone-create" // bool
	rolesCreatePermissionToggleKey     = "cacophony:processor:cc:roles-create"    // roleIDs comma delimted
	denyServerCommandsUsePermissionKey = "cacophony:processor:cc:server-use"      // bool (true means can't use)
	denyUserCommandsUsePermissionKey   = "cacophony:processor:cc:user-use"        // bool (true means can't use)
	editQuestionnaireKey               = "cacophony:processor:dev:edit-questionnaire"
	deleteQuestionnaireKey             = "cacophony:processor:dev:remove-questionnaire"
)

func seporateUserAndServerCommands(commands []CustomCommand) (userCommands []CustomCommand, serverCommands []CustomCommand) {
	for _, entry := range commands {
		if entry.IsUserCommand {
			userCommands = append(userCommands, entry)
		} else {
			serverCommands = append(serverCommands, entry)
		}
	}
	return
}
