package ping

import (
	"time"

	"github.com/bwmarrin/discordgo"

	"gitlab.com/Cacophony/go-kit/events"
)

func handlePing(event *events.Event) {
	createdAt, err := event.MessageCreate.Timestamp.Parse()
	if err != nil {
		event.Except(err)
		return
	}

	_, err = event.RespondComplex(
		&discordgo.MessageSend{
			Embed: &discordgo.MessageEmbed{
				Title: "ping.ping-response.embed.title",
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "ping.ping-response.embed.field.DiscordToGateway.title",
						Value:  "ping.ping-response.embed.field.DiscordToGateway.value",
						Inline: true,
					},
					{
						Name:   "ping.ping-response.embed.field.GatewayToProcessor.title",
						Value:  "ping.ping-response.embed.field.GatewayToProcessor.value",
						Inline: true,
					},
				},
			},
		},
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

func handlePang(event *events.Event) {
	_, err := event.Respond("ping.pang-response")
	if err != nil {
		event.Except(err)
	}
}
