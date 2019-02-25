package dev

import (
	"strings"

	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handleBotOwners(event *events.Event) {
	message := "The following users are Bot Owners: "
	for _, botOwnerID := range event.BotOwnerIDs() {
		user, err := p.state.User(botOwnerID)
		if err != nil {
			message += "`#" + botOwnerID + "`"
			continue
		}

		message += "`" + user.String() + "`"
		message += ", "
	}
	message = strings.TrimRight(message, ", ")

	_, err := event.Respond(message)
	event.Except(err)
}
