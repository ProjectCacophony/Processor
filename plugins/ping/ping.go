package ping

import (
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
	"gitlab.com/Cacophony/go-kit/discord"

	"gitlab.com/Cacophony/go-kit/events"
)

func (p *Plugin) handlePing(event *events.Event) {
	createdAt, err := event.MessageCreate.Timestamp.Parse()
	if err != nil {
		event.Except(err)
		return
	}

	discordToGateway := event.ReceivedAt.Sub(createdAt).Round(time.Millisecond)
	gatewayToProcessor := time.Since(event.ReceivedAt).Round(time.Millisecond)

	var proxyDuration time.Duration
	// only test proxy ping if a proxy is configured
	if discordgo.EndpointDiscord != "https://discordapp.com/" &&
		discordgo.EndpointDiscord != "https://discord.com/" {
		beforeProxy := time.Now()
		proxyResp, err := event.HTTPClient().Get(discordgo.EndpointDiscord + "status")
		proxyDuration = time.Since(beforeProxy).Round(time.Millisecond)
		if err != nil || proxyResp.StatusCode != http.StatusOK {
			event.Except(err)
			return
		}
	}

	sendStart := time.Now()
	messages, err := event.RespondComplex(
		&discordgo.MessageSend{
			Embed: pingEmbed(),
		},
		"DiscordToGateway",
		discordToGateway,
		"GatewayToProcessor",
		gatewayToProcessor,
		"ProxyLatency",
		proxyDuration,
	)
	if err != nil {
		event.Except(err)
		return
	}
	sendDuration := time.Since(sendStart).Round(time.Millisecond)

	if len(messages) <= 0 {
		return
	}

	_, err = discord.EditComplexWithVars(
		nil,
		event.Discord(),
		event.Localizations(),
		&discordgo.MessageEdit{
			Embed:   pingEmbed(),
			ID:      messages[0].ID,
			Channel: messages[0].ChannelID,
		},
		false,
		"DiscordToGateway",
		discordToGateway,
		"GatewayToProcessor",
		gatewayToProcessor,
		"ProxyLatency",
		proxyDuration,
		"SendDuration",
		sendDuration,
	)
	event.Except(err)
}

func (p *Plugin) handlePong(event *events.Event) {
	_, err := event.Respond("ping.pong-response")
	if err != nil {
		event.Except(err)
	}
}

func (p *Plugin) handlePang(event *events.Event) {
	_, err := event.Respond("ping.pang-response")
	if err != nil {
		event.Except(err)
	}
}

func pingEmbed() *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: "ping.ping-response.embed.title",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "ping.ping-response.embed.field.Event.title",
				Value:  "ping.ping-response.embed.field.Event.value",
				Inline: false,
			},
			{
				Name:   "ping.ping-response.embed.field.SendDuration.title",
				Value:  "ping.ping-response.embed.field.SendDuration.value",
				Inline: false,
			},
		},
	}
}
