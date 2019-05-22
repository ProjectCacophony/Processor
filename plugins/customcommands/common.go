package customcommands

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
