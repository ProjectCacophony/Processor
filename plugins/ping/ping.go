package ping

import (
	"time"

	"gitlab.com/Cacophony/go-kit/events"
)

func handlePing(event *events.Event) {
	createdAt, err := event.MessageCreate.Timestamp.Parse()
	if err != nil {
		event.Except(err)
		return
	}

	_, err = event.Respond(
		"ping.ping-response",
		"DiscordToGateway",
		event.ReceivedAt.Sub(createdAt),
		"GatewayToProcessor",
		time.Since(event.ReceivedAt),
	)
	if err != nil {
		event.Except(err)
		return
	}
}

func handlePong(event *events.Event) {
	_, err := event.Respond("ping.pong-response")
	if err != nil {
		event.Except(err)
	}
}
