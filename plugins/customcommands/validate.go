package customcommands

import (
	"strings"

	"gitlab.com/Cacophony/go-kit/events"
)

func isUserOperation(e *events.Event) bool {
	return e.DM() || hasUserParam(e)
}

func hasUserParam(e *events.Event) bool {
	return (len(e.Fields()) >= 3 && e.Fields()[2] == "user")
}

func isValidCommandName(name string) bool {
	return !strings.Contains(name, " ")
}
