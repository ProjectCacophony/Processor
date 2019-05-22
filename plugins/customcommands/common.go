package customcommands

const (
	everyoneCreatePermissionKey    = "cacophony:processor:cc:everyone-create" // bool
	rolesCreatePermissionToggleKey = "cacophony:processor:cc:roles-create"    // roleIDs comma delimted
	serverCommandsUsePermissionKey = "cacophony:processor:cc:server-use"      // bool (true means can't use)
	userCommandsUsePermissionKey   = "cacophony:processor:cc:user-use"        // bool (true means can't use)
)

func seporateUserAndServerEntries(entries []Entry) (userEntries []Entry, serverEntries []Entry) {
	for _, entry := range entries {
		if entry.IsUserCommand {
			userEntries = append(userEntries, entry)
		} else {
			serverEntries = append(serverEntries, entry)
		}
	}
	return
}
