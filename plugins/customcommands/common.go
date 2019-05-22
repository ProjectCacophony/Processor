package customcommands

const (
	everyoneCreatePermissionKey    = "cacophony:processor:prefix:everyone-create" // true/false
	rolesCreatePermissionToggleKey = "cacophony:processor:prefix:roles-create"    // roleIDs comma delimted
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
